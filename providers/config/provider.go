// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package config

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
	"os"
)

const Name = "config"

type ConfigProvider struct {
	container.ServiceProvider
	mainPath string // 二进制 main 程序的绝对路径
}

func (self *ConfigProvider) Name() string {
	return Name
}

func (self *ConfigProvider) Boot(c container.Container) error {
	file, err := os.Getwd()
	if err == nil {
		self.mainPath = file + "/"
	} else {
		logger.Instance().Trace(err.Error())
	}
	logger.Instance().Trace("Boot Path Provider.")
	return nil
}

func (self *ConfigProvider) RegisterProviderInstance(c container.Container) container.NewInstance {
	return func(params ...interface{}) (interface{}, error) {
		c := params[0].(container.Container)
		return &ConfigService{container: c, mainPath: self.mainPath}, nil
	}
}

func (*ConfigProvider) IsDefer() bool {
	return false
}

func (sp *ConfigProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}
