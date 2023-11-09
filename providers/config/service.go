package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/kit/castkit"
	"github.com/text3cn/goodle/kit/filekit"
	"github.com/text3cn/goodle/kit/strkit"
	"github.com/text3cn/goodle/types"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Instance() *ConfigService {
	if instance == nil {
		file, _ := os.Getwd()
		instance = &ConfigService{currentPath: file + "/"}
	}
	return instance
}

type ConfigService struct {
	Service
	container   core.Container
	currentPath string       // 二进制 main 程序的绝对路径
	lock        sync.RWMutex // 配置文件读写锁
}

type Service interface {
	LoadConfig(filename string) (*viper.Viper, error)
	Get(key string) *castkit.GoodleVal

	IsDebug() bool

	GetHttpAddr() string
	GetRuntimePath() string
	GetDatabase() (dbsCfg map[string]types.DBConfig)
	GetRedis() map[string]types.RedisConfig
	GetSwagger() types.SwaggerConfig
	GetEtcd() types.EtcdConfig
	GetGoodLog() types.Goodlog
}

func (self *ConfigService) LoadConfig(filename string) (*viper.Viper, error) {
	//goodlog.Pink("self.currentPath: ", self.currentPath)
	if configs[filename] != nil {
		return configs[filename], nil
	}
	var config, localConfig, mergerConfig *viper.Viper
	seg := strings.Split(filename, ".")
	fName := seg[0]
	fType := seg[1]
	cfgFile := fName + "." + fType
	// 在当前目录找配置文件
	if exists, _ := filekit.PathExists(filepath.Join(self.currentPath, cfgFile)); exists {
		config = loadConfigFile(self.currentPath, fName, fType)
	}
	// 当前目录没找到，到 ./config 目录中找
	if config == nil {
		path := filepath.Join(self.currentPath, "config")
		if exists, _ := filekit.PathExists(path); exists {
			config = loadConfigFile(path, fName, fType)
		}
	}
	// 本地配置文件
	localConfigDir := filepath.Join(filepath.Dir(self.currentPath), "config", "local")
	localConfigFile := filepath.Join(localConfigDir, cfgFile)
	if exists, _ := filekit.PathExists(localConfigFile); exists {
		localConfig = loadConfigFile(localConfigDir, fName, fType)
	}
	// 没有配置文件
	if config == nil && localConfig == nil {
		err := errors.New("Unable to find configuration file " + filename +
			". The configuration file should be placed in any of the following paths:\n " +
			self.currentPath + filename + "\n" +
			self.currentPath + "config/" + filename + "\n" +
			self.currentPath + "config/local/" + filename + "\n",
		)
		return nil, err
	}
	// 如有本地配置文件，则用本地配置项覆盖
	mergerConfig = config
	if localConfig != nil {
		allKeys := localConfig.AllKeys()
		for _, v := range allKeys {
			mergerConfig.Set(v, localConfig.Get(v))
		}
	}
	// 缓存起来
	configs[fName] = mergerConfig
	return mergerConfig, nil
}

func (self *ConfigService) Get(key string) *castkit.GoodleVal {
	seg := strkit.Explode(".", key)
	if len(seg) == 1 {
		return &castkit.GoodleVal{}
	}
	cfg := configs[seg[0]]
	itemKey := strings.Replace(key, seg[0]+".", "", 1)
	return &castkit.GoodleVal{cfg.Get(itemKey)}
}

func (self *ConfigService) IsDebug() bool {
	key := "debug"
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
	return self.currentPath
}

func (self *ConfigService) GetDatabase() (dbsCfg map[string]types.DBConfig) {
	dbsCfg = make(map[string]types.DBConfig)
	cfg, _ := self.LoadConfig("database.yaml")
	cfgNodes := mergerLevel2(cfg)
	for k, v := range cfgNodes {
		item := types.DBConfig{}
		mapstructure.Decode(v, &item)
		dbsCfg[k] = item
	}
	return
}

func (self *ConfigService) GetRedis() (configs map[string]types.RedisConfig) {
	configs = make(map[string]types.RedisConfig)
	cfg, _ := self.LoadConfig("redis.yaml")
	if cfg != nil {
		cfgNodes := mergerLevel2(cfg)
		for k, v := range cfgNodes {
			item := types.RedisConfig{}
			mapstructure.Decode(v, &item)
			configs[k] = item
		}
	}
	return
}

func (self *ConfigService) GetSwagger() (config types.SwaggerConfig) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("swagger")
		mapstructure.Decode(value, &config)
	}
	return
}

func (self *ConfigService) GetEtcd() (config types.EtcdConfig) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("etcd")
		mapstructure.Decode(value, &config)
	}
	return
}

func (self *ConfigService) GetGoodLog() (config types.Goodlog) {
	if cfg, _ := self.getAppConfig(); cfg != nil {
		value := cfg.Get("goodlog")
		mapstructure.Decode(value, &config)
	}
	return
}
