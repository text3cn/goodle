package orm

import (
	"gorm.io/gorm"
)

type builder struct {
	db         *gorm.DB
	sql        string
	bindValues []any
}

func Builder(sql string, bindValues ...any) builder {
	return builder{
		sql: sql,
		db:  GetDB(),
	}
}
