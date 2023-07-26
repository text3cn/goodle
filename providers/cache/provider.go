// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package cache

import (
	 "github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
)

const Name = "cache"

type CacheServiceProvider struct {
	container.ServiceProvider
}

func (self *CacheServiceProvider) Name() string {
	return Name
}

func (sp *CacheServiceProvider) RegisterProviderInstance(c container.Container) container.NewInstance {
	return func(params ...interface{}) (interface{}, error) {

		c := params[0].(container.Container)
		return &CacheService{c: c}, nil
	}
}

func (*CacheServiceProvider) IsDefer() bool {
	return false
}

func (sp *CacheServiceProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (sp *CacheServiceProvider) Boot(c container.Container) error {
	logger.Instance().Trace("Boot Cache Provider")
	return nil
}
