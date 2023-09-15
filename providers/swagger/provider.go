package swagger

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
)

const Name = "swagger"

var instance *SwaggerService

type SwaggerProvider struct {
	container.ServiceProvider
}

func (self *SwaggerProvider) Name() string {
	return Name
}

func (*SwaggerProvider) InitOnBoot() bool {
	return true
}

func (*SwaggerProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (*SwaggerProvider) RegisterProviderInstance(c container.Container) container.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &SwaggerService{c: c}
		return instance, nil
	}
}

func (*SwaggerProvider) BeforeInit(c container.Container) error {
	logger.Instance().Trace("BeforeInit Swagger Provider")
	return nil
}

func (*SwaggerProvider) AfterInit(instance any) error {
	cfg := instance.(Service)
	logger.Instance().Trace("AfterInit Swagger Provider")
	cfg.Init()
	return nil
}
