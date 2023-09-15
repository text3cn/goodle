package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/kit/castkit"
	"os"
	"sync"
)

var instance *ConfigService

func Instance() *ConfigService {
	if instance == nil {
		file, _ := os.Getwd()
		instance = &ConfigService{mainPath: file + "/"}
	}
	return instance
}

var allConfigItem = map[string]interface{}{}

// 使用泛型会有编译开销和执行效率开销，也许这点开销对于业务而言并不关心，
// 但是作为一个框架 goodle 的理念是尽量不要去增加开销，所以为每个配置项定义了方法
type Service interface {
	Get(key string) *castkit.GoodleVal
	GetHttpAddr() string
	IsDevelop() bool
	IsDaemon() bool
	GetRuntimePath() string
	LoadConfig(filename string) (*viper.Viper, error)
	GetDatabase() (dbsCfg map[string]DBConfig)
	GetRedis() map[string]RedisConfig
	GetSwagger() SwaggerConfig
}

// path 服务建立在 config 服务之上，框架需用使用路径的时候直接找 path 服务拿
// path 服务封装了取 config 和取默认路径的细节
type ConfigService struct {
	Service
	container container.Container
	mainPath  string       // 二进制 main 程序的绝对路径
	lock      sync.RWMutex // 配置文件读写锁
}

func (self *ConfigService) Get(key string) *castkit.GoodleVal {
	return &castkit.GoodleVal{allConfigItem[key]}
}

func (self *ConfigService) getAppConfig() (*viper.Viper, error) {
	cfg, err := self.LoadConfig("app.yaml")
	return cfg, err
}

// http server listen config
func (self *ConfigService) GetHttpAddr() string {
	addr := ""
	if cfg, _ := self.getAppConfig(); cfg != nil {
		if value, ok := cfg.Get("http.port").(int); !ok {
			panic("The configuration of http.port is not a valid value")
		} else {
			addr = ":" + cast.ToString(value)
		}
	}
	return addr
}

// 是否后台运行
func (self *ConfigService) IsDevelop() bool {
	key := "develop"
	if cfg, _ := self.getAppConfig(); cfg != nil {
		if cfg.IsSet(key) {
			if val, ok := cfg.Get(key).(bool); !ok {
				panic("The configuration of " + key + " is not a valid value")
			} else {
				return cast.ToBool(val)
			}
		}
	}
	return false
}

// 是否以守护进程方式运行
func (self *ConfigService) IsDaemon() bool {
	key := "daemon"
	if cfg, _ := self.getAppConfig(); cfg != nil {
		if cfg.IsSet(key) {
			if val, ok := cfg.Get(key).(bool); !ok {
				panic("The configuration of " + key + " is not a valid value")
			} else {
				return cast.ToBool(val)
			}
		}
	}
	return false
}

// runtime 目录
func (self *ConfigService) GetRuntimePath() string {
	key := "runtime.path"
	if cfg, _ := self.getAppConfig(); cfg != nil {
		if cfg.IsSet(key) {
			if val, ok := cfg.Get(key).(string); !ok {
				panic("The configuration of " + key + " is not a valid value")
			} else {
				return cast.ToString(val)
			}
		}
	}
	return self.mainPath
}

func (self *ConfigService) GetDatabase() (dbsCfg map[string]DBConfig) {
	dbsCfg = make(map[string]DBConfig)
	cfg, _ := self.LoadConfig("database.yaml")
	cfgNodes := mergerLevel2(cfg)
	for k, v := range cfgNodes {
		item := DBConfig{}
		mapstructure.Decode(v, &item)
		dbsCfg[k] = item
	}
	return
}

func (self *ConfigService) GetRedis() (configs map[string]RedisConfig) {
	configs = make(map[string]RedisConfig)
	cfg, _ := self.LoadConfig("redis.yaml")
	if cfg != nil {
		cfgNodes := mergerLevel2(cfg)
		for k, v := range cfgNodes {
			item := RedisConfig{}
			mapstructure.Decode(v, &item)
			configs[k] = item
		}
	}
	return
}

func (self *ConfigService) GetSwagger() (config SwaggerConfig) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("swagger")
		mapstructure.Decode(value, &config)
	}
	return
}

func (self *ConfigService) GetDiscovery() (config EtcdConfig) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("etcd")
		mapstructure.Decode(value, &config)
	}
	return
}
