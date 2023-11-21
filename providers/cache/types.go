package cache

import (
	"github.com/redis/go-redis/v9"
)

type SourceSetter func() (cache string)

type redisConfig struct {
	conn          *redis.Client
	name          string // conn name
	expireSeconds int
}

type freeCachePipelineSet struct {
	Cache        string
	write        bool
	cacheHolder  *FreeCacheHolder
	key          string
	widthLocal   bool
	localExpire  int
	redisClients []redisConfig
}

type fastCachePipelineSet struct {
	Data         any
	Cache        string
	write        bool
	cacheHolder  *FastCacheHolder
	key          string
	widthLocal   bool
	redisClients []redisConfig
}
