package cache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"github.com/text3cn/goodle/providers/goodlog"
	goodleredis "github.com/text3cn/goodle/providers/redis"
	"time"
)

var bigCacheInstance *bigCache

type bigCache struct {
	buckets map[string]*BigCacheHolder
}

type BigCacheHolder struct {
	*bigcache.BigCache
}

// 开辟自动扩容桶
func NewBigCache(bucketName string, expireSeconds int) {
	if bigCacheInstance == nil {
		bigCacheInstance = &bigCache{buckets: make(map[string]*BigCacheHolder)}
	}
	var cache *bigcache.BigCache
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(time.Duration(expireSeconds)*time.Second))
	bigCacheInstance.buckets[bucketName] = &BigCacheHolder{cache}
}

// 自定义配置开辟桶
func NewBigCacheWithConfig(bucketName string, config bigcache.Config) {
	cache, initErr := bigcache.New(context.Background(), config)
	if initErr != nil {
		goodlog.Error(initErr)
		return
	}
	bigCacheInstance.buckets[bucketName] = &BigCacheHolder{cache}
}

// 获取桶
func (s *CacheService) BigCache(bucketName string) *BigCacheHolder {
	return bigCacheInstance.buckets[bucketName]
}

// ////////////////////////////////// 扩展方法 ////////////////////////////////////
func (self *BigCacheHolder) Set(key, value string) {
	self.BigCache.Set(key, []byte(value))
}

func (self *BigCacheHolder) Get(key string) string {
	got, err := self.BigCache.Get(key)
	if err != nil {
		return ""
	}
	return string(got)
}

func (self *BigCacheHolder) Delete(key string) bool {
	err := self.BigCache.Delete(key)
	self.BigCache.Len()
	if err != nil {
		goodlog.Error(err.Error())
		return false
	}
	return true
}

func (self *BigCacheHolder) Clear() {
	self.BigCache.Reset()
}

//////////////////////////////////// Pipeline ////////////////////////////////////

func (self *BigCacheHolder) Pipeline(key string, write ...bool) *bigCachePipelineSet {
	w := false
	if len(write) > 0 {
		w = true
	}
	return &bigCachePipelineSet{key: key, cacheHolder: self, write: w}
}

func (self *bigCachePipelineSet) Local() *bigCachePipelineSet {
	self.widthLocal = true
	if self.Data != "" {
		return self
	}
	goodlog.Trace("BigCache 中查找 " + self.key)
	cache := self.cacheHolder.Get(self.key)
	if cache != "" {
		self.Data = cache
	}
	return self
}

// args[0]: expire seconds
// args[1]: connection name
func (self *bigCachePipelineSet) Redis(args ...any) *bigCachePipelineSet {
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
	self.cacheHolder.Set(self.key, value)
	return self
}

func (self *bigCachePipelineSet) Setter(setter SourceSetter) *bigCachePipelineSet {
	if self.Data != "" && self.write == false {
		return self
	}
	goodlog.Trace("Setter 中产生数据 ", self.key)
	cache := setter()
	if cache == "" {
		return self
	}
	self.Data = cache
	// 回写缓存
	if self.widthLocal {
		self.cacheHolder.Set(self.key, cache)
	}
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			expire := time.Duration(redis.expireSeconds) * time.Second
			redis.conn.Set(context.Background(), self.key, cache, expire)
		}
	}
	return self
}

func (self *bigCachePipelineSet) Delete() {
	self.Data = ""
	self.cacheHolder.Delete(self.key)
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			redis.conn.Del(context.Background(), self.key)
		}
	}
}

func (self *bigCachePipelineSet) Clear() {
	self.Data = ""
	self.cacheHolder.Delete(self.key)
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			redis.conn.Del(context.Background(), self.key)
		}
	}
}
