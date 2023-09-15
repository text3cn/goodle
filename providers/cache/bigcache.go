package cache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/text3cn/goodle/providers/logger"
	"time"
)

var bigCacheInstance *bigCache

type bigCache struct {
	buckets map[string]*bigCacheHolder
}

type bigCacheHolder struct {
	*bigcache.BigCache
}

// 开辟自动扩容桶
func NewBigCache(bucketName string, expireSeconds int) {
	if bigCacheInstance == nil {
		bigCacheInstance = &bigCache{buckets: make(map[string]*bigCacheHolder)}
	}
	var cache *bigcache.BigCache
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(time.Duration(expireSeconds)*time.Second))
	bigCacheInstance.buckets[bucketName] = &bigCacheHolder{cache}
}

// 自定义配置开辟桶
func NewBigCacheWithConfig(bucketName string, config bigcache.Config) {
	cache, initErr := bigcache.New(context.Background(), config)
	if initErr != nil {
		logger.Instance().Error(initErr)
		return
	}
	bigCacheInstance.buckets[bucketName] = &bigCacheHolder{cache}
}

// 获取桶
func (s *CacheService) BigCache(bucketName string) *bigCacheHolder {
	return bigCacheInstance.buckets[bucketName]
}

// ////////////////////////////////// 扩展方法 ////////////////////////////////////
func (self *bigCacheHolder) Set(key, value string) {
	self.BigCache.Set(key, []byte(value))
}

func (self *bigCacheHolder) Get(key string) string {
	got, err := self.BigCache.Get(key)
	if err != nil {
		return ""
	}
	return string(got)
}

func (self *bigCacheHolder) Delete(key string) bool {
	err := self.BigCache.Delete(key)
	self.BigCache.Len()
	if err != nil {
		logger.Instance().Error(err.Error())
		return false
	}
	return true
}

func (self *bigCacheHolder) Clear() {
	self.BigCache.Reset()
}
