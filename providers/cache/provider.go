// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package cache

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/goodlog"
)

const Name = "cache"

var instance *CacheService

type CacheServiceProvider struct {
	core.ServiceProvider
}

func (self *CacheServiceProvider) Name() string {
	return Name
}

func (sp *CacheServiceProvider) RegisterProviderInstance(holder core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &CacheService{holder: holder}
		return instance, nil
	}
}

func (*CacheServiceProvider) InitOnBoot() bool {
	return false
}

func (sp *CacheServiceProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (sp *CacheServiceProvider) BeforeInit(c core.Container) error {
	goodlog.Trace("BeforeInit Cache Provider")
	return nil
}

func (*CacheServiceProvider) AfterInit(instance any) error {
	return nil
}
