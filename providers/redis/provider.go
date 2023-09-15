package redis

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
	"sync"
)

const Name = "redis"

var instance *RedisService

//func Instance(c container.ServicesContainer, name ...string) *redis.Client {
//	if instance == nil {
//		redisService := c.NewSingle(Name).(Service)
//		redisService.init()
//	}
//	key := "first"
//	if len(name) > 0 {
//		key = name[0]
//	}
//	return instance.dbs[key]
//}

type ReidsProvider struct {
	container.ServiceProvider
}

func (self *ReidsProvider) Name() string {
	return Name
}

func (*ReidsProvider) RegisterProviderInstance(c container.Container) container.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &RedisService{c: c, lock: sync.Mutex{}}
		return instance, nil
	}
}

func (*ReidsProvider) InitOnBoot() bool {
	return false
}

func (*ReidsProvider) BeforeInit(c container.Container) error {
	logger.Instance().Trace("BeforeInit Reids Provider")
	return nil
}

func (*ReidsProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (*ReidsProvider) AfterInit(instance any) error {
	return nil
}
