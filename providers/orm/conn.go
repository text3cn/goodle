package orm

import (
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

var isDebug = false

// 在 GetDB 时进行服务初始化，连接数据库
// 使用 map 保存多个数据，dsn 作为 key 确保单例
func GetDB(connName ...string) *gorm.DB {
	if instance == nil {
		ormService := core.FrameContainer().NewSingle(Name).(Service)
		ormService.Init()
	}
	var key string
	cfg := config.Instance().GetDatabase()
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

func mysqlOpen(config types.DBConfig) gorm.Dialector {
	isDebug = config.Debug
	return mysql.New(mysql.Config{
		DSN:                       formatDsn(config),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})
}

// 生成 dsn
// https://gorm.io/zh_CN/docs/connecting_to_the_database.html
func formatDsn(conf types.DBConfig) (dsn string) {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn += conf.Username + ":" + conf.Password
	dsn += "@" + conf.Protocol + "(" + conf.Host + ":" + strconv.Itoa(conf.Port) + ")"
	dsn += "/" + conf.Database
	dsn += "?charset=" + conf.Charset + "&parseTime=" + cast.ToString(conf.ParseTime) + "&loc=" + conf.Loc
	return
}
