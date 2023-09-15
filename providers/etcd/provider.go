package etcd

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
)

const Name = "etcd"

var etcdInstance *etcdService

type ReidsProvider struct {
	container.ServiceProvider
}

func (self *ReidsProvider) Name() string {
	return Name
}

func (*ReidsProvider) RegisterProviderInstance(box container.Container) container.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		etcdInstance = &etcdService{c: box}
		return etcdInstance, nil
	}
}

func (*ReidsProvider) InitOnBoot() bool {
	return false
}

func (*ReidsProvider) BeforeInit(c container.Container) error {
	logger.Instance().Trace("BeforeInit Etcd Provider")
	return nil
}

func (*ReidsProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (*ReidsProvider) AfterInit(instance any) error {
	return nil
}
