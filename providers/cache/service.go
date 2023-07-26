package cache

import (
	"github.com/text3cn/goodle/container"
)

type Service interface {
	LocalCache(output string)
}

type CacheService struct {
	Service
	c container.Container
}

func (s *CacheService) LocalCache(out string) {
	println("[缓存]", out)
}