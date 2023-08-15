package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/kit/castkit"
	"github.com/text3cn/goodle/kit/filekit"
	"github.com/text3cn/goodle/providers/logger"
	"github.com/text3cn/goodle/types"
	"path/filepath"
	"strings"
	"sync"
)

var allConfigItem = map[string]interface{}{}

// 使用泛型会有编译开销和执行效率开销，也许这点开销对于业务而言并不关心，
// 但是作为一个框架 goodle 的理念是尽量不要去增加开销，所以为每个配置项定义了方法
type Service interface {
	Get(key string) *castkit.GoodleVal
	GetHttpAddr() string
	IsDevelop() bool
	GetRuntimePath() string
	LoadConfig(filename string) *viper.Viper
	GetDatabase() (dbsCfg map[string]types.DBConfig)
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

// 优先使用可执行程序当前目录下的 xxx.yaml，
// 然后找 config 目录下的 xxx.yaml，如果有则将两者合并，前者覆盖后者
func (self *ConfigService) LoadConfig(filename string) *viper.Viper {
	//logger.Pink("self.mainPath: ", self.mainPath)
	noneModuleConfig := false
	var cmdConfig, commonConfig, mergerConfig *viper.Viper
	seg := strings.Split(filename, ".")
	fName := seg[0]
	fType := seg[1]
	cfgFile := fName + "." + fType
	// 模块 config
	cmdCfgFile := filepath.Join(self.mainPath, cfgFile)
	if exists, _ := filekit.PathExists(cmdCfgFile); exists {
		cfg := viper.New()
		cfg.AddConfigPath(self.mainPath)
		cfg.SetConfigName(fName)
		cfg.SetConfigType(fType)
		if err := cfg.ReadInConfig(); err != nil {
			panic(err)
		}
		cmdConfig = cfg
	} else {
		noneModuleConfig = true
	}
	// 公共 config
	configDir := filepath.Dir(filepath.Dir(filepath.Dir(self.mainPath)))
	configDir = filepath.Join(configDir, "config")
	defaultConfigFile := filepath.Join(configDir, cfgFile)
	if exists, _ := filekit.PathExists(defaultConfigFile); exists {
		cfg := viper.New()
		cfg.AddConfigPath(configDir)
		cfg.SetConfigName(fName)
		cfg.SetConfigType(fType)
		if err := cfg.ReadInConfig(); err != nil {
			panic(err)
		}
		commonConfig = cfg
	} else {
		if noneModuleConfig {
			logger.Pink("Unable to find configuration file. The configuration file needs to be placed in any of the following directories:")
			fmt.Println(self.mainPath + filename)
			fmt.Println(defaultConfigFile + "\n")
		}
	}
	// merger
	if commonConfig != nil {
		mergerConfig = commonConfig
	}
	if cmdConfig != nil {
		allKeys := cmdConfig.AllKeys()
		for _, v := range allKeys {
			mergerConfig.Set(v, cmdConfig.Get(v))
		}
	}
	// 存一份暴露给用户使用
	for _, v := range mergerConfig.AllKeys() {
		if allConfigItem[v] == nil {
			allConfigItem[v] = mergerConfig.Get(v)
		}
	}
	return mergerConfig
}

func (self *ConfigService) getDefaultConfig() *viper.Viper {
	return self.LoadConfig("app.yaml")
}

// http server listen config
func (self *ConfigService) GetHttpAddr() string {
	addr := ":80"
	if cfg := self.getDefaultConfig(); cfg != nil {
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
	if cfg := self.getDefaultConfig(); cfg != nil {
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
	if cfg := self.getDefaultConfig(); cfg != nil {
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

func (self *ConfigService) GetDatabase() (dbsCfg map[string]types.DBConfig) {
	dbsCfg = make(map[string]types.DBConfig)
	cfg := self.LoadConfig("database.yaml")
	cfgNodes := mergerLevel2(cfg)
	for k, v := range cfgNodes {
		item := types.DBConfig{}
		mapstructure.Decode(v, &item)
		dbsCfg[k] = item
	}
	return
}

// 将多个二级配置合并公共的一级配置项，相同配置项用二级覆盖一级
func mergerLevel2(source *viper.Viper) (ret map[string]map[string]interface{}) {
	ret = make(map[string]map[string]interface{})
	allkeys := source.AllKeys()
	common := map[string]interface{}{}
	nodes := map[string]map[string]interface{}{}
	for _, v := range allkeys {
		seg := strings.Split(v, ".")
		if len(seg) == 1 {
			common[v] = source.Get(v)
		} else {
			if nodes[seg[0]] == nil {
				item := make(map[string]interface{})
				item[seg[1]] = source.Get(v)
				nodes[seg[0]] = item
			} else {
				nodes[seg[0]][seg[1]] = source.Get(v)
			}
		}
	}
	for key, node := range nodes {
		it := map[string]interface{}{}
		for k, v := range common {
			it[k] = v
		}
		for k, v := range node {
			it[k] = v
		}
		ret[key] = make(map[string]interface{})
		ret[key] = it
	}
	return
}
