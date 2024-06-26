package orm

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/goodlog"
)

const Name = "orm"

var instance *OrmService
var log goodlog.Service
var configService config.Service

type OrmProvider struct {
	core.ServiceProvider
}

func (self *OrmProvider) Name() string {
	return Name
}

func (*OrmProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &OrmService{c: c}
		return instance, nil
	}
}

func (*OrmProvider) InitOnBoot() bool {
	return false
}

func (*OrmProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*OrmProvider) BeforeInit(c core.Container) error {
	log.Trace("BeforeInit Orm Provider")
	return nil
}

func (*OrmProvider) AfterInit(instance any) error {
	return nil
}
