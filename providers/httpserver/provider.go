// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package httpserver

import (
	container "github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
)

const Name = "httpserver"

type HttpServerProvider struct {
	container.ServiceProvider
	HttpServer *Engine
}

func (self *HttpServerProvider) Name() string {
	return Name
}

func (self *HttpServerProvider) BeforeInit(c container.Container) error {
	logger.Instance().Trace("BeforeInit HttpServer Provider")
	//self.HttpServer
	return nil
}

func (sp *HttpServerProvider) RegisterProviderInstance(c container.Container) container.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		c := params[0].(container.Container)
		return &HttpServerService{container: c}, nil
	}
}

func (*HttpServerProvider) InitOnBoot() bool {
	return true
}

func (sp *HttpServerProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (*HttpServerProvider) AfterInit(instance any) error {
	return nil
}
