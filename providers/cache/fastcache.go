package cache

import (
	"github.com/VictoriaMetrics/fastcache"
)

var fastCacheInstance *fastCache

type fastCache struct {
	buckets map[string]*fastCacheHolder
}

type fastCacheHolder struct {
	*fastcache.Cache
}

// 开辟自动扩容桶
func NewFastCache(bucketName string, size int) {
	if fastCacheInstance == nil {
		fastCacheInstance = &fastCache{buckets: make(map[string]*fastCacheHolder)}
	}
	// 容量，字节为单位，小于 32 MB 当做 32 MB 处理
	var cache = fastcache.New(size)
	fastCacheInstance.buckets[bucketName] = &fastCacheHolder{cache}
}

// 获取桶
func (s *CacheService) FastCache(bucketName string) *fastCacheHolder {
	return fastCacheInstance.buckets[bucketName]
}

// ////////////////////////////////// 扩展方法 ////////////////////////////////////
func (self *fastCacheHolder) Set(key, value string) {
	self.Cache.Set([]byte(key), []byte(value))
}

func (self *fastCacheHolder) Get(key string) string {
	got := self.Cache.Get(nil, []byte(key))
	return string(got)
}

func (self *fastCacheHolder) Delete(key string) {
	self.Cache.Del([]byte(key))
}

func (self *fastCacheHolder) Clear() {
	self.Cache.Reset()
}
