package orm

func CheckEmptyId(sql string, bindValues ...any) bool {
	id := 0
	GetDB().Raw(sql+" LIMIT 1", bindValues...).Scan(&id)
	if id > 0 {
		return true
	}
	return false
}
