package orm

import (
	"database/sql"
	goodleconfig "github.com/text3cn/goodle/config"
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Service interface {
	Init()
	//GetDB(option ...DBOption) (*gorm.DB, error)
}

type OrmService struct {
	Service
	c   container.Container
	dbs map[string]*gorm.DB // key 为 dsn, value 为 gorm.DB（连接池）

}

// 如果能获取到配置文件则进行数据库连接，这个目前在 engine/run.go -> initServices() 中直接实例化了
func (self *OrmService) Init() {
	configService := goodleconfig.Instance()
	dbsCfg := configService.GetDatabase()
	if dbsCfg == nil {
		logger.Pink("database config error")
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
		// 设置对应的连接池配置
		sqlDB, err = db.DB()
		if err != nil {
			break
		}
		if config.ConnMaxIdle > 0 {
			sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
		}
		if config.ConnMaxOpen > 0 {
			sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
		}
		if config.ConnMaxLifetime != "" {
			liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
			if err != nil {
				logger.Pink("conn max lift time error", map[string]interface{}{"err": err})
			} else {
				sqlDB.SetConnMaxLifetime(liftTime)
			}
		}
		if config.ConnMaxIdletime != "" {
			idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
			if err != nil {
				logger.Pink("conn max idle time error", map[string]interface{}{"err": err})
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
