package orm

import (
	"github.com/text3cn/goodle/core"
	"gorm.io/gorm"
)

type builder struct {
	db         *gorm.DB
	sql        string
	bindValues []any
}

func Builder(c core.Container, sql string, bindValues ...any) builder {
	return builder{
		sql: sql,
		db:  GetDB(c),
	}
}
