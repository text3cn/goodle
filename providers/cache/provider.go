// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package cache

import (
	servicescenter2 "github.com/text3cn/t3web/container"
	"github.com/text3cn/t3web/providers/logger"
)

const Name = "cache"

type CacheServiceProvider struct {
	servicescenter2.ServiceProvider
}

func (self *CacheServiceProvider) Name() string {
	return Name
}

func (sp *CacheServiceProvider) RegisterProviderInstance(c servicescenter2.Container) servicescenter2.NewInstance {
	return func(params ...interface{}) (interface{}, error) {

		c := params[0].(servicescenter2.Container)
		return &CacheService{c: c}, nil
	}
}

func (*CacheServiceProvider) IsDefer() bool {
	return false
}

func (sp *CacheServiceProvider) Params(c servicescenter2.Container) []interface{} {
	return []interface{}{c}
}

func (sp *CacheServiceProvider) Boot(c servicescenter2.Container) error {
	logger.Instance().Trace("Boot Cache Provider")
	return nil
}
