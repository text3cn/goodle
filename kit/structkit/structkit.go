package structkit

import "reflect"

// 将结构体 src 中与 dst 相同的字段值扫描到结构体 dst 中
// 用法：structkit.CopyStruct(item, &it)
func CopyStruct(src, dst interface{}) {
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst).Elem()
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		fieldName := srcValue.Type().Field(i).Name
		dstField := dstValue.FieldByName(fieldName)
		if dstField.IsValid() && dstField.CanSet() {
			dstField.Set(srcField)
		}
	}
}
