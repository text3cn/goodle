package cache

import (
	"context"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"github.com/text3cn/goodle/providers/goodlog"
	goodleredis "github.com/text3cn/goodle/providers/redis"
	"time"
)

var freeCacheInstance *freeCache

type freeCache struct {
	buckets map[string]*FreeCacheHolder
}

type FreeCacheHolder struct {
	*freecache.Cache
}

// 开辟桶
func NewFreeCache(bucketName string, size int) {
	if freeCacheInstance == nil {
		freeCacheInstance = &freeCache{buckets: make(map[string]*FreeCacheHolder)}
	}
	freeCacheInstance.buckets[bucketName] = &FreeCacheHolder{freecache.NewCache(size)}
}

// 获取桶
func (s *CacheService) FreeCache(bucketName string) *FreeCacheHolder {
	return freeCacheInstance.buckets[bucketName]
}

// ////////////////////////////////// 扩展方法 ////////////////////////////////////
func (self *FreeCacheHolder) Set(key, value string, expireSeconds int) {
	k := []byte(key)
	v := []byte(value)
	self.Cache.Set(k, v, expireSeconds)
}

func (self *FreeCacheHolder) Get(key string) string {
	k := []byte(key)
	got, err := self.Cache.Get(k)
	if err != nil {
		return ""
	}
	return string(got)
}

func (self *FreeCacheHolder) Delete(key string) bool {
	k := []byte(key)
	return self.Cache.Del(k)
}

func (self *FreeCacheHolder) Clear() {
	self.Cache.Clear()
}

//////////////////////////////////// Pipeline ////////////////////////////////////

func (self *FreeCacheHolder) Pipeline(key string, write ...bool) *freeCachePipelineSet {
	w := false
	if len(write) > 0 {
		w = true
	}
	return &freeCachePipelineSet{key: key, cacheHolder: self, write: w}
}

func (self *freeCachePipelineSet) Local(expireSeconds ...int) *freeCachePipelineSet {
	self.localExpire = 0
	if len(expireSeconds) > 0 {
		self.localExpire = expireSeconds[0]
	}
	self.widthLocal = true
	if self.Data != "" {
		return self
	}
	goodlog.Trace("\nFastCache 中查找 " + self.key)
	cache := self.cacheHolder.Get(self.key)
	if cache != "" {
		self.Data = cache
	}
	return self
}

// args[0]: expire seconds
// args[1]: connection name
func (self *freeCachePipelineSet) Redis(args ...any) *freeCachePipelineSet {
	redisSvc := instance.holder.NewSingle(goodleredis.Name).(goodleredis.Service)
	var client *redis.Client
	expire := 0
	connName := ""
	argsLen := len(args)
	if argsLen == 0 {
		client = redisSvc.Conn()
	} else if len(args) == 1 {
		if e, ok := args[0].(int); ok {
			expire = e
		}
		client = redisSvc.Conn()
	} else if len(args) == 2 {
		if c, ok := args[1].(string); ok {
			connName = c
			client = redisSvc.Conn(connName)
		}
	}
	self.redisClients = append(self.redisClients, redisConfig{
		conn:          client,
		name:          connName,
		expireSeconds: expire,
	})
	if self.Data != "" {
		return self
	}
	goodlog.Trace("Redis 中查找 " + self.key)
	cache := client.Get(context.Background(), self.key)
	value := cache.Val()
	if value != "" {
		self.Data = value
	}
	// 回写到本地缓存
	self.cacheHolder.Set(self.key, value, self.localExpire)
	return self
}

func (self *freeCachePipelineSet) Setter(setter SourceSetter) *freeCachePipelineSet {
	if self.Data != "" && self.write == false {
		return self
	}
	goodlog.Error("Setter 中产生数据 ", self.key)
	cache := setter()
	if cache == "" {
		return self
	}
	self.Data = cache
	// 回写缓存
	if self.widthLocal {
		self.cacheHolder.Set(self.key, cache, self.localExpire)
	}
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			expire := time.Duration(redis.expireSeconds) * time.Second
			redis.conn.Set(context.Background(), self.key, cache, expire)

		}
	}
	return self
}

func (self *freeCachePipelineSet) Delete() {
	self.Data = ""
	self.cacheHolder.Delete(self.key)
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			redis.conn.Del(context.Background(), self.key)
		}
	}
}

func (self *freeCachePipelineSet) Clear() {
	self.Data = ""
	self.cacheHolder.Delete(self.key)
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			redis.conn.Del(context.Background(), self.key)
		}
	}
}
