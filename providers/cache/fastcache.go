package cache

import (
	"github.com/VictoriaMetrics/fastcache"
)

var fastCacheInstance *fastCache

type fastCache struct {
	buckets map[string]*FastCacheHolder
}

type FastCacheHolder struct {
	*fastcache.Cache
}

// 开辟自动扩容桶
func NewFastCache(bucketName string, size int) *FastCacheHolder {
	if fastCacheInstance == nil {
		fastCacheInstance = &fastCache{buckets: make(map[string]*FastCacheHolder)}
	}
	// 容量，字节为单位，小于 32 MB 当做 32 MB 处理
	var cache = fastcache.New(size)
	fastCacheInstance.buckets[bucketName] = &FastCacheHolder{cache}
	return fastCacheInstance.buckets[bucketName]
}

// 获取桶
func (s *CacheService) FastCache(bucketName string) *FastCacheHolder {
	return fastCacheInstance.buckets[bucketName]
}

// ////////////////////////////////// 扩展方法 ////////////////////////////////////
func (self *FastCacheHolder) Set(key, value string) {
	self.Cache.Set([]byte(key), []byte(value))
}

func (self *FastCacheHolder) Get(key string) string {
	got := self.Cache.Get(nil, []byte(key))
	return string(got)
}

func (self *FastCacheHolder) Delete(key string) {
	self.Cache.Del([]byte(key))
}

func (self *FastCacheHolder) Clear() {
	self.Cache.Reset()
}

/* fastcache 没有过期时间，所以没有支持 Pipeline */
