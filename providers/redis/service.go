package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/text3cn/goodle/config"
	"github.com/text3cn/goodle/container"
	"sync"
)

type Service interface {
	init()
	Conn(connName ...string) *redis.Client
	Pipeline(string, ...bool) *pipeline
}

type RedisService struct {
	Service
	c    container.Container
	dbs  map[string]*redis.Client
	lock sync.Mutex
}

// 如果能获取到配置文件则进行连接
func (self *RedisService) init() {
	self.lock.Lock()
	defer self.lock.Unlock()
	configService := config.Instance()
	redisConfigs := configService.GetRedis()
	var isFirst = true
	for connName, config := range redisConfigs {
		rdb := redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.Db,
		})
		// 挂到 map 中
		if self.dbs == nil {
			self.dbs = make(map[string]*redis.Client)
		}
		if isFirst {
			connName = "first"
			isFirst = false
		}
		self.dbs[connName] = rdb
	}
}

func (self *RedisService) Conn(connName ...string) *redis.Client {
	if self.dbs == nil {
		self.init()
	}
	name := "first"
	if len(connName) > 0 {
		name = connName[0]
	}
	return self.dbs[name]
}
