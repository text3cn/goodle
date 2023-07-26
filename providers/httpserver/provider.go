// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package httpserver

import (
	servicescenter2 "github.com/text3cn/t3web/container"
	"github.com/text3cn/t3web/providers/logger"
)

const Name = "httpserver"

type HttpServerProvider struct {
	servicescenter2.ServiceProvider
	HttpServer *T3WebEngine
}

func (self *HttpServerProvider) Name() string {
	return Name
}

func (self *HttpServerProvider) Boot(c servicescenter2.Container) error {
	logger.Instance().Trace("Boot HttpServer Provider")
	//self.HttpServer
	return nil
}

func (sp *HttpServerProvider) RegisterProviderInstance(c servicescenter2.Container) servicescenter2.NewInstance {
	return func(params ...interface{}) (interface{}, error) {
		c := params[0].(servicescenter2.Container)
		return &HttpServerService{container: c}, nil
	}
}

func (*HttpServerProvider) IsDefer() bool {
	return false
}

func (sp *HttpServerProvider) Params(c servicescenter2.Container) []interface{} {
	return []interface{}{c}
}
