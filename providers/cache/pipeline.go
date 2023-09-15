package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/text3cn/goodle/providers/logger"
	goodleredis "github.com/text3cn/goodle/providers/redis"
	"time"
)

func (self *freeCacheHolder) Pipeline(key string, write ...bool) *freeCachePipelineSet {
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
	if self.Cache != "" {
		return self
	}
	logger.Instance().Trace("\nFastCache 中查找 " + self.key)
	cache := self.cacheHolder.Get(self.key)
	if cache != "" {
		self.Cache = cache
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
	if self.Cache != "" {
		return self
	}
	logger.Instance().Trace("Redis 中查找 " + self.key)
	cache := client.Get(context.Background(), self.key)
	value := cache.Val()
	if value != "" {
		self.Cache = value
	}
	return self
}

func (self *freeCachePipelineSet) Setter(setter SourceSetter) *freeCachePipelineSet {
	if self.Cache != "" && self.write == false {
		return self
	}
	logger.Pink("Setter 中产生数据 ", self.key)
	cache := setter()
	self.Cache = cache
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
	self.Cache = ""
	self.cacheHolder.Delete(self.key)
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			redis.conn.Del(context.Background(), self.key)
		}
	}
}

func (self *freeCachePipelineSet) Clear() {
	self.Cache = ""
	self.cacheHolder.Delete(self.key)
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			redis.conn.Del(context.Background(), self.key)
		}
	}
}
