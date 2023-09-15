package orm

import (
	"github.com/text3cn/goodle/container"
	"gorm.io/gorm"
)

type builder struct {
	db         *gorm.DB
	sql        string
	bindValues []any
}

func Builder(c container.Container, sql string, bindValues ...any) builder {
	return builder{
		sql: sql,
		db:  GetDB(c),
	}
}
