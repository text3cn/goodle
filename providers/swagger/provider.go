package swagger

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/goodlog"
)

const Name = "swagger"

var instance *SwaggerService

type SwaggerProvider struct {
	core.ServiceProvider
}

func (self *SwaggerProvider) Name() string {
	return Name
}

func (*SwaggerProvider) InitOnBoot() bool {
	return true
}

func (*SwaggerProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*SwaggerProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &SwaggerService{c: c}
		return instance, nil
	}
}

func (*SwaggerProvider) BeforeInit(c core.Container) error {
	goodlog.Trace("BeforeInit Swagger Provider")
	return nil
}

func (*SwaggerProvider) AfterInit(instance any) error {
	cfg := instance.(Service)
	goodlog.Trace("AfterInit Swagger Provider")
	cfg.Init()
	return nil
}
