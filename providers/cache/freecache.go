package cache

import (
	"github.com/coocood/freecache"
)

var freeCacheInstance *freeCache

type freeCache struct {
	buckets map[string]*FreeCacheHolder
}

type FreeCacheHolder struct {
	*freecache.Cache
}

// 开辟桶
func NewFreeCache(bucketName string, size int) {
	if freeCacheInstance == nil {
		freeCacheInstance = &freeCache{buckets: make(map[string]*FreeCacheHolder)}
	}
	freeCacheInstance.buckets[bucketName] = &FreeCacheHolder{freecache.NewCache(size)}
}

// 获取桶
func (s *CacheService) FreeCache(bucketName string) *FreeCacheHolder {
	return freeCacheInstance.buckets[bucketName]
}

// ////////////////////////////////// 扩展方法 ////////////////////////////////////
func (self *FreeCacheHolder) Set(key, value string, expireSeconds int) {
	k := []byte(key)
	v := []byte(value)
	self.Cache.Set(k, v, expireSeconds)
}

func (self *FreeCacheHolder) Get(key string) string {
	k := []byte(key)
	got, err := self.Cache.Get(k)
	if err != nil {
		return ""
	}
	return string(got)
}

func (self *FreeCacheHolder) Delete(key string) bool {
	k := []byte(key)
	return self.Cache.Del(k)
}

func (self *FreeCacheHolder) Clear() {
	self.Cache.Clear()
}
