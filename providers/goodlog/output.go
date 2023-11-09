package goodlog

import (
	"fmt"
	"github.com/spf13/cast"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func (self *GoodlogService) output(level byte, v ...interface{}) error {
	// calldepth 层次说明：本函数（第一层） -> service.go（第二层） -> level.go（第三层）-> 用户调用时（第四层）
	// 最后一个参数可以传 deep:number 来设置层数，例如：deep:5
	length := len(v)
	calldepth := 4
	last := cast.ToString(v[length-1])
	if len(last) >= 4 && last[:4] == "deep" {
		arr := strings.Split(last, ":")
		calldepth, _ = strconv.Atoi(arr[1])
		length--
	}
	// 格式化输出样式
	str := ""
	formatStr := "%v"
	for i := 0; i < length; i++ {
		str += fmt.Sprintf(formatStr, v[i]) + " "
	}
	switch level {
	case trace:
		formatStr = "\033[37m[TRACE] %s\033[0m"
	case debug:
		formatStr = "\033[32m[DEBUG] %s\033[0m"
	case info:
		formatStr = "\033[36m[INFO] %s\033[0m"
	case warn:
		formatStr = "\033[33m[WARN] %s\033[0m"
	case err:
		formatStr = "\033[35m[ERROR] %s\033[0m"
	case fatal:
		formatStr = "\033[31m[FATAL] %s\033[0m"
	}
	str = fmt.Sprintf(formatStr, str)
	//	s := fmt.Sprintf(formatStr, v...)
	// fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, 1, 97, 31, message, 0x1B)
	//s = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 1, 97, 31, v, 0x1B)

	_log := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
	return _log.Output(calldepth, str)
}

func (s *GoodlogService) P(color string, data any) {
	_log := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
	calldepth := 4
	// 类型
	rVal := reflect.ValueOf(data)
	rKind := rVal.Kind().String()
	rType := rVal.Type().String()
	typeStr := fmt.Sprintf(yellowBG+black+" %v -> %v "+reset, rKind, rType)

	// 正文
	fmt.Print(color)
	switch rKind {
	case "array", "slice":
		switch rType {
		case "[]int":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]int)
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", v)
			}
		case "[]int64":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]int64)
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", v)
			}
		case "[]string":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]string)
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", v)
			}
		case "[]uint8":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			fmt.Printf("%v", string(rVal.Bytes()))
		case "[]map[string]interface {}":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]map[string]interface{})
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", "map (")
				for key, val := range v {
					fmt.Println("      ", "[", key, "]", "=>", val)
				}
				fmt.Println("),")
			}
		case "[]interface {}":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]interface{})
			for k, v := range array {
				fmt.Println("\t", "[", k, "]", "=>", v)
			}
		default:
			str := fmt.Sprintf(color+" %v "+reset, data)
			_log.Output(4, typeStr+str)
			fmt.Print(color)
		}
	case "map":
		switch rType {
		case "map[string]interface {}":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			for k, v := range data.(map[string]interface{}) {
				fmt.Println("\t", "[", k, "]", "=>", v)
			}
		default:
			str := fmt.Sprintf(color+" %v "+reset, data)
			_log.Output(4, typeStr+str)
			fmt.Print(color)
		}
	case "struct":
		_log.Output(calldepth, typeStr)
		fmt.Print(color)
		t := reflect.TypeOf(data)
		v := reflect.ValueOf(data)
		for k := 0; k < t.NumField(); k++ {
			fmt.Println("\t", "[", t.Field(k).Name, "]", "=>", v.Field(k).Interface())
		}
	default:
		str := fmt.Sprintf(color+" %v "+reset, data)
		_log.Output(4, typeStr+str)
		fmt.Print(color)
	}
	// 重置样式
	fmt.Printf(reset)
}
