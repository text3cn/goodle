package orm

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/goodlog"
	"gorm.io/gorm"
)

var isDebug = false

var container *core.ServicesContainer

// 在 GetDB 时进行服务初始化，连接数据库
// 使用 map 保存多个数据，dsn 作为 key 确保单例
func GetDB(connName ...string) *gorm.DB {
	initServices()
	var key string
	cfg := configService.GetDatabase()
	for _, v := range cfg {
		key = v.DefaultConn
		break
	}
	if len(connName) > 0 {
		key = connName[0]
	}
	if isDebug {
		return instance.dbs[key].Debug()
	}
	return instance.dbs[key]
}

func initServices() {
	if container == nil {
		container = core.NewContainer()
		container.Bind(&config.ConfigProvider{})
		container.Bind(&goodlog.GoodlogProvider{})
		container.Bind(&OrmProvider{})
		configService = container.NewSingle(config.Name).(config.Service)
		log = container.NewSingle(goodlog.Name).(goodlog.Service)
	}
	if instance == nil {
		ormService := container.NewSingle(Name).(Service)
		ormService.Init()
	}
}
