// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package path

import (
	servicescenter2 "github.com/text3cn/t3web/container"
	"github.com/text3cn/t3web/providers/logger"
	"os"
)

const Name = "path"

type PathProvider struct {
	servicescenter2.ServiceProvider
	RootPath string
}

func (self *PathProvider) Name() string {
	return Name
}

func (self *PathProvider) Boot(c servicescenter2.Container) error {
	logger.Instance().Trace("Boot Path Provider")
	file, err := os.Getwd()
	if err == nil {
		self.RootPath = file + "/"
	}
	return nil
}

func (sp *PathProvider) RegisterProviderInstance(c servicescenter2.Container) servicescenter2.NewInstance {
	return func(params ...interface{}) (interface{}, error) {
		c := params[0].(servicescenter2.Container)
		return &PathService{container: c}, nil
	}
}

func (*PathProvider) IsDefer() bool {
	return false
}

func (sp *PathProvider) Params(c servicescenter2.Container) []interface{} {
	return []interface{}{c}
}
