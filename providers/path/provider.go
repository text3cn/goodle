// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package path

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
	"os"
)

const Name = "path"

type PathProvider struct {
	container.ServiceProvider
	RootPath string
}

func (self *PathProvider) Name() string {
	return Name
}

func (self *PathProvider) Boot(c container.Container) error {
	logger.Instance().Trace("Boot Path Provider")
	file, err := os.Getwd()
	if err == nil {
		self.RootPath = file + "/"
	}
	return nil
}

func (sp *PathProvider) RegisterProviderInstance(c container.Container) container.NewInstance {
	return func(params ...interface{}) (interface{}, error) {
		c := params[0].(container.Container)
		return &PathService{container: c}, nil
	}
}

func (*PathProvider) IsDefer() bool {
	return false
}

func (sp *PathProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}
