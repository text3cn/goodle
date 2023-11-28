package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/providers/goodlog"
	"time"
)

type pipeline struct {
	write bool // 是否强制重写缓存，即刷新缓存
	//redisSvc *RedisService
	key string
}

type pipelineSet struct {
	Data         string
	write        bool
	key          string
	redisClients []redisConfig
}

type SourceSetter func() (cache string)

type redisConfig struct {
	conn          *redis.Client
	name          string // conn name
	expireSeconds int
}

func (self *RedisService) Pipeline(key string, write ...bool) *pipeline {
	w := false
	if len(write) > 0 {
		w = true
	}
	return &pipeline{key: key, write: w}
}

// args[0]: expire seconds
// args[1]: connection name
func (self *pipeline) Redis(args ...any) *pipelineSet {
	ret := &pipelineSet{write: self.write, key: self.key}
	var client *redis.Client
	expire := 0
	connName := ""
	argsLen := len(args)
	if argsLen == 0 {
		client = instance.Conn()
	} else if len(args) == 1 {
		if e, ok := args[0].(int); ok {
			expire = e
		} else {
			goodlog.Error("Expire sencods " + cast.ToString(args[0]) + " is not a valid value.")
			return ret
		}
		client = instance.Conn()
	} else if len(args) == 2 {
		if c, ok := args[1].(string); ok {
			connName = c
			client = instance.Conn(connName)
		}
	}
	ret.redisClients = append(ret.redisClients, redisConfig{
		conn:          client,
		name:          connName,
		expireSeconds: expire,
	})
	if ret.Data != "" {
		return ret
	}
	goodlog.Trace("Redis " + connName + " 中查找 " + self.key)
	cache := client.Get(context.Background(), self.key)
	value := cache.Val()
	if value != "" {
		ret.Data = value
	}
	return ret
}

func (self *pipelineSet) Setter(setter SourceSetter) *pipelineSet {
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
	if self.redisClients != nil {
		for _, redis := range self.redisClients {
			expire := time.Duration(redis.expireSeconds) * time.Second
			redis.conn.Set(context.Background(), self.key, cache, expire)

		}
	}
	return self
}
