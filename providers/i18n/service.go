package i18n

import (
	"github.com/spf13/viper"
	"github.com/text3cn/goodle/core"
)

//var pkgs = make(map[string]*viper.Viper)

type Service interface {
	SetLngCode(string)
	LoadedPackage(string) bool
	SetLanguagePackage(string, *viper.Viper)
	I18n() func(string) string
}

type I18nService struct {
	Service // 实现接口，显示标记
	c       core.Container
	lngCode string
	pkgs    map[string]*viper.Viper
}

func (self *I18nService) SetLngCode(lng string) {
	self.lngCode = lng
}

// 判断语言包是否已加载到内存
func (self *I18nService) LoadedPackage(lng string) bool {
	return self.pkgs[lng] != nil
}

// 从磁盘加载到内存
func (self *I18nService) SetLanguagePackage(lng string, pkg *viper.Viper) {
	self.pkgs[lng] = pkg
}

func (self *I18nService) I18n() func(string) string {
	return func(key string) string {
		return self.pkgs[self.lngCode].GetString(key)
	}
}
