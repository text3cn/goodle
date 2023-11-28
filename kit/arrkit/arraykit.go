package arrkit

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// 从 string 数组中删除一个指定值，返回删除后的新数组
func ArrayStringDeleteVal(val string, array []string) []string {
	// 找出值的 index
	index := 0
	count := len(array)
	for i := 0; i < count; i++ {
		if val == array[i] {
			index = i
			break
		}
	}
	// 删除
	array = append(array[:index], array[index+1:]...)
	return array
}

// 从数组中删除一个元素
// index 要删除的元素下标
func ArrayStringDeleteOne(index int, array []string) []string {
	array = append(array[:index], array[index+1:]...)
	return array
}

// 判断符合数据类型中是否存在某个值
// 支持的类型: slice 、array 、map
func InArray(value interface{}, array interface{}) bool {
	val := reflect.ValueOf(array)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(value, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(value, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: 只支持类型为 slice、array、map 的数据验证")
	}
	return false
}

// Int 数组倒序
func IntArrayDesc(arr []int) (ret []int) {
	sort.Ints(arr)
	for i := len(arr) - 1; i >= 0; i-- {
		ret = append(ret, arr[i])
	}
	return ret
}

// 数组转逗号拼接
func JoinWithCommas(numbers []int) string {
	var strNumbers []string
	for _, num := range numbers {
		strNumbers = append(strNumbers, fmt.Sprint(num))
	}
	return strings.Join(strNumbers, ",")
}
