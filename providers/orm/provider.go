package orm

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
)

const Name = "orm"

var instance *OrmService

type OrmProvider struct {
	container.ServiceProvider
}

func (self *OrmProvider) Name() string {
	return Name
}

func (*OrmProvider) RegisterProviderInstance(c container.Container) container.NewInstance {
	return func(params ...interface{}) (interface{}, error) {
		instance = &OrmService{c: c}
		return instance, nil
	}
}

func (*OrmProvider) IsDefer() bool {
	return false
}

func (*OrmProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (*OrmProvider) Boot(c container.Container) error {
	logger.Pink("Boot Orm Provider")
	return nil
}
