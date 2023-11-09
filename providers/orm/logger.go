package orm

import (
	"context"
	"github.com/text3cn/goodle/providers/goodlog"
	"time"
)

// orm 的日志实现类, 实现了gorm.Logger.Interface
type OrmLogger struct {
	logger *goodlog.GoodlogService // 用一个logger对象存放hade的log服务
}

// NewOrmLogger 初始化一个ormLogger,
//func NewOrmLogger() *OrmLogger {
//	return &OrmLogger{logger: logger.Instance()}
//}

func (o *OrmLogger) Info(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	//o.logger.Info(ctx, s, fields)
	o.logger.Info(s, fields)
}

func (o *OrmLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	//o.logger.Warn(ctx, s, fields)
	o.logger.Warn(s, fields)
}

func (o *OrmLogger) Error(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	o.logger.Error(ctx, s, fields)
}

func (o *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	fields := map[string]interface{}{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}
	s := "orm trace sql"
	o.logger.Trace(ctx, s, fields)
}
