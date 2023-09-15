package orm

import "github.com/text3cn/goodle/container"

func CheckEmptyId(c container.Container, sql string, bindValues ...any) bool {
	id := 0
	GetDB(c).Raw(sql+" LIMIT 1", bindValues...).Scan(&id)
	if id > 0 {
		return true
	}
	return false
}
