package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/text3cn/goodle/kit/filekit"
	"path/filepath"
	"strings"
)

// 默认使用可执行程序当前目录下的 xxx.yaml，
// 如果 developtime 目录存在配置文件则覆盖配置，主要用于本地开发，developtime 中的配置不要提交到版本库
func (self *ConfigService) LoadConfig(filename string) (*viper.Viper, error) {
	//logger.Pink("self.mainPath: ", self.mainPath)
	noneReleaseConfig := false
	var releaseConfig, developConfig, mergerConfig *viper.Viper
	seg := strings.Split(filename, ".")
	fName := seg[0]
	fType := seg[1]
	cfgFile := fName + "." + fType
	// 模块可执行文件所在目录的 config
	releaseCfgFile := filepath.Join(self.mainPath, cfgFile)
	if exists, _ := filekit.PathExists(releaseCfgFile); exists {
		cfg := viper.New()
		cfg.AddConfigPath(self.mainPath)
		cfg.SetConfigName(fName)
		cfg.SetConfigType(fType)
		if err := cfg.ReadInConfig(); err != nil {
			panic(err)
		}
		releaseConfig = cfg
	} else {
		noneReleaseConfig = true
	}
	// developtime config
	configDir := filepath.Dir(self.mainPath)
	configDir = filepath.Join(configDir, "developtime")
	developConfigFile := filepath.Join(configDir, cfgFile)
	if exists, _ := filekit.PathExists(developConfigFile); exists {
		cfg := viper.New()
		cfg.AddConfigPath(configDir)
		cfg.SetConfigName(fName)
		cfg.SetConfigType(fType)
		if err := cfg.ReadInConfig(); err != nil {
			panic(err)
		}
		developConfig = cfg
	} else {
		if noneReleaseConfig {
			err := errors.New("Unable to find configuration file " + filename +
				". The configuration file needs to be placed in " + self.mainPath + filename)
			return nil, err
		}
	}
	// merger，用开发配置覆盖发行配置
	mergerConfig = releaseConfig
	if developConfig != nil {
		allKeys := developConfig.AllKeys()
		for _, v := range allKeys {
			mergerConfig.Set(v, developConfig.Get(v))
		}
	}
	// 存一份暴露给用户使用
	for _, v := range mergerConfig.AllKeys() {
		if allConfigItem[v] == nil {
			allConfigItem[v] = mergerConfig.Get(v)
		}
	}
	return mergerConfig, nil
}

// 在同一个配置问价那种，将多个二级配置合并公共的一级配置项，相同配置项用二级覆盖一级，这种配置文件写法主要是省去重复配置
// 例如一个配置文件中连接多个数据库时提取公共配置为一级，每个库的差异化配置做为二级。
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
