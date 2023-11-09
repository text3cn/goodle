// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package httpserver

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/goodlog"
)

const Name = "httpserver"

type HttpServerProvider struct {
	core.ServiceProvider
	HttpServer *Engine
}

func (self *HttpServerProvider) Name() string {
	return Name
}

func (self *HttpServerProvider) BeforeInit(c core.Container) error {
	goodlog.Trace("BeforeInit HttpServer Provider")
	//self.HttpServer
	return nil
}

func (sp *HttpServerProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		c := params[0].(core.Container)
		return &HttpServerService{container: c}, nil
	}
}

func (*HttpServerProvider) InitOnBoot() bool {
	return true
}

func (sp *HttpServerProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*HttpServerProvider) AfterInit(instance any) error {
	return nil
}
