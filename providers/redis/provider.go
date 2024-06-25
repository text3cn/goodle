package redis

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/goodlog"
	"sync"
)

const Name = "redis"

var instance *RedisService

func Instance() *RedisService {
	if instance == nil {
		redisService := core.NewContainer().NewSingle(Name).(*RedisService)
		redisService.init()
		return redisService
	}
	return instance
}

type ReidsProvider struct {
	core.ServiceProvider
}

func (self *ReidsProvider) Name() string {
	return Name
}

func (*ReidsProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &RedisService{c: c, lock: sync.Mutex{}}
		return instance, nil
	}
}

func (*ReidsProvider) InitOnBoot() bool {
	return false
}

func (*ReidsProvider) BeforeInit(c core.Container) error {
	goodlog.Trace("BeforeInit Reids Provider")
	return nil
}

func (*ReidsProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*ReidsProvider) AfterInit(instance any) error {
	return nil
}
