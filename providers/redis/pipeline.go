package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/providers/goodlog"
)

type pipeline struct {
	Cache        string
	write        bool
	redisSvc     *RedisService
	key          string
	redisClients []redisConfig
}

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
	return &pipeline{key: key, redisSvc: self, write: w}
}

// args[0]: expire seconds
// args[1]: connection name
func (self *pipeline) Redis(args ...any) *pipeline {
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
			return self
		}
		client = instance.Conn()
	} else if len(args) == 2 {
		if c, ok := args[1].(string); ok {
			connName = c
			client = instance.Conn(connName)
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
	goodlog.Trace("Redis " + connName + " 中查找 " + self.key)
	cache := client.Get(context.Background(), self.key)
	value := cache.Val()
	if value != "" {
		self.Cache = value
	}
	return self
}
