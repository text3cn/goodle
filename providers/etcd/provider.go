package etcd

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/goodlog"
)

const Name = "etcd"

var etcdInstance *etcdService

type ReidsProvider struct {
	core.ServiceProvider
}

func (self *ReidsProvider) Name() string {
	return Name
}

func (*ReidsProvider) RegisterProviderInstance(box core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		etcdInstance = &etcdService{c: box}
		return etcdInstance, nil
	}
}

func (*ReidsProvider) InitOnBoot() bool {
	return false
}

func (*ReidsProvider) BeforeInit(c core.Container) error {
	goodlog.Trace("BeforeInit Etcd Provider")
	return nil
}

func (*ReidsProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*ReidsProvider) AfterInit(instance any) error {
	return nil
}
