// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package cache

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
)

const Name = "cache"

var instance *CacheService

type CacheServiceProvider struct {
	container.ServiceProvider
}

func (self *CacheServiceProvider) Name() string {
	return Name
}

func (sp *CacheServiceProvider) RegisterProviderInstance(holder container.Container) container.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &CacheService{holder: holder}
		return instance, nil
	}
}

func (*CacheServiceProvider) InitOnBoot() bool {
	return false
}

func (sp *CacheServiceProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (sp *CacheServiceProvider) BeforeInit(c container.Container) error {
	logger.Instance().Trace("BeforeInit Cache Provider")
	return nil
}

func (*CacheServiceProvider) AfterInit(instance any) error {
	return nil
}
