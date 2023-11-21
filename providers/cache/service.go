package cache

import (
	"github.com/text3cn/goodle/core"
)

type Service interface {
	FreeCache(string) *FreeCacheHolder
	BigCache(string) *BigCacheHolder
	FastCache(string) *FastCacheHolder
}

type CacheService struct {
	Service
	holder core.Container
}
