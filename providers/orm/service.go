package orm

import (
	"database/sql"
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/goodlog"
	"github.com/text3cn/goodle/types"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

type Service interface {
	Init()
	//GetDB(option ...DBOption) (*gorm.DB, error)
}

type OrmService struct {
	Service
	c   core.Container
	dbs map[string]*gorm.DB // key 为 dsn, value 为 gorm.DB（连接池）

}

// 如果能获取到配置文件则进行数据库连接
func (self *OrmService) Init() {
	dbsCfg := configService.GetDatabase()
	if dbsCfg == nil {
		goodlog.Error("database config error")
		return
	}
	for k, config := range dbsCfg {
		// 实例化 gorm.DB
		var db *gorm.DB
		var err error
		var sqlDB *sql.DB
		switch config.Driver {
		case "mysql":
			db, err = gorm.Open(mysqlOpen(config), &gorm.Config{
				// 禁用自动迁移创建的表名称复数形式
				NamingStrategy: schema.NamingStrategy{SingularTable: true},
			})
		case "sqlite":
			db, err = gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
		}
		// 设置对应的连接池配置，确保 instance 中随时可以拿到未断开的连接
		sqlDB, err = db.DB()
		if err != nil {
			break
		}
		connMaxIdle := 10
		maxOpenConns := 100
		if config.ConnMaxIdle > 0 {
			connMaxIdle = config.ConnMaxIdle
		}
		if config.ConnMaxOpen > 0 {
			maxOpenConns = config.ConnMaxOpen
		}
		sqlDB.SetMaxIdleConns(connMaxIdle)
		sqlDB.SetMaxOpenConns(maxOpenConns)

		if config.ConnMaxLifetime != "" {
			liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
			if err != nil {
				goodlog.Error("conn max lift time error", map[string]interface{}{"err": err})
			} else {
				sqlDB.SetConnMaxLifetime(liftTime)
			}
		}
		if config.ConnMaxIdletime != "" {
			idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
			if err != nil {
				goodlog.Error("conn max idle time error", map[string]interface{}{"err": err})
			} else {
				sqlDB.SetConnMaxIdleTime(idleTime)
			}
		}
		// 挂载到 map 中
		if self.dbs == nil {
			self.dbs = make(map[string]*gorm.DB)
		}
		self.dbs[k] = db
	}
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
