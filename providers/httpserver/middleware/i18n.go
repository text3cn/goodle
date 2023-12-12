package middleware

import (
	"github.com/spf13/viper"
	"github.com/text3cn/goodle/kit/filekit"
	"github.com/text3cn/goodle/providers/goodlog"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/i18n"
	"path/filepath"
)

// dir 语言包文件所在目录（绝对路径）
func I18n(dir string) httpserver.MiddlewareHandler {
	return func(ctx *httpserver.Context) error {
		lngCode := ctx.GetVal("lng").ToString()
		languagePkg := dir + string(filepath.Separator) + lngCode + ".json"
		goodlog.Trace("Use i18n middleware. dir:", languagePkg)
		i18nService := ctx.Holder().NewSingle(i18n.Name).(i18n.Service)
		i18nService.SetLngCode(lngCode)
		if !i18nService.LoadedPackage(lngCode) {
			if exists, _ := filekit.PathExists(languagePkg); exists {
				file := viper.New()
				file.AddConfigPath(dir)
				file.SetConfigName(lngCode + ".json")
				file.SetConfigType("json")
				if err := file.ReadInConfig(); err != nil {
					goodlog.Error(err)
				} else {
					i18nService.SetLanguagePackage(lngCode, file)
				}
			} else {
				goodlog.Error("Cant find language package:", languagePkg)
			}
		}
		ctx.Next()
		return nil
	}
}
