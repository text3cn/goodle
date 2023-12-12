package config

import (
	"github.com/spf13/viper"
	"strings"
)

var configs = make(map[string]*viper.Viper)

func loadConfigFile(dir, fName, fType string) *viper.Viper {
	cfg := viper.New()
	cfg.AddConfigPath(dir)
	cfg.SetConfigName(fName)
	cfg.SetConfigType(fType)
	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}
	return cfg
}

// 在同一个配置文件中，将多个二级配置合并公共的一级配置项，相同配置项用二级覆盖一级，这种配置文件写法主要是省去重复配置
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

func (self *ConfigService) getAppConfig() (*viper.Viper, error) {
	cfg, err := self.LoadConfig("app.yaml")
	return cfg, err
}
