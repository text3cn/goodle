package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/kit/filekit"
	"path/filepath"
)

// 使用泛型会有编译开销和执行效率开销，也许这点开销对于业务而言并不关心，
// 但是作为一个框架 goodle 的理念是尽量不要去增加开销，所以为每个配置项定义了方法
type Service interface {
	LogPath()
	GetHttpAddr() string
	IsDevelop() bool
	GetRuntimePath() string
}

// path 服务建立在 config 服务之上，框架需用使用路径的时候直接找 path 服务拿
// path 服务封装了取 config 和取默认路径的细节
type ConfigService struct {
	Service
	container container.Container
	mainPath  string // 二进制 main 程序的绝对路径
	// 一些其他路径，比如日志输出路径等，通过配置文件传递进来，配置文件目录等
}

// 优先使用可执行程序当前目录下的 app.yaml，
// 然后找 config 目录下的 app.yaml，如果有则将两者合并，前者覆盖后者
func (self *ConfigService) getDefaultConfig() *viper.Viper {
	//fmt.Println("self.mainPath", self.mainPath)
	var cmdConfig, commonConfig, mergerConfig *viper.Viper
	// cmd config
	cmdCfgFile := filepath.Join(self.mainPath, "app.yaml")
	if exists, _ := filekit.PathExists(cmdCfgFile); exists {
		cfg := viper.New()
		cfg.AddConfigPath(self.mainPath)
		cfg.SetConfigName("app")
		cfg.SetConfigType("yaml")
		if err := cfg.ReadInConfig(); err != nil {
			panic(err)
		}
		cmdConfig = cfg
	}
	// config
	configDir := filepath.Dir(filepath.Dir(filepath.Dir(self.mainPath)))
	configDir = filepath.Join(configDir, "config")
	defaultConfigFile := filepath.Join(configDir, "app.yaml")
	if exists, _ := filekit.PathExists(defaultConfigFile); exists {
		cfg := viper.New()
		cfg.AddConfigPath(configDir)
		cfg.SetConfigName("app")
		cfg.SetConfigType("yaml")
		if err := cfg.ReadInConfig(); err != nil {
			panic(err)
		}
		commonConfig = cfg
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
	return mergerConfig
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

func (s *ConfigService) LogPath() {

}
