package cache

import (
	"github.com/text3cn/goodle/core"
)

type Service interface {
	FreeCache(string) *freeCacheHolder
	BigCache(string) *bigCacheHolder
	FastCache(string) *fastCacheHolder
}

type CacheService struct {
	Service
	holder core.Container
}
