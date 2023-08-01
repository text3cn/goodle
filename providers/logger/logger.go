package logger

import (
	"fmt"
	"github.com/spf13/cast"
)

func Pink(output ...interface{}) {
	outputWithColor(35, output...)
}

func outputWithColor(color int, v ...interface{}) {
	// 最后一个参数可以传 deep:数字来设置 calldepth
	length := len(v)
	// 格式化输出样式
	str := ""
	for i := 0; i < length; i++ {
		str += cast.ToString(v[i])
	}
	str = fmt.Sprintf("\n%c[%d;%d;%dm%s%c[0m", 0x1B, 1, 97, color, str, 0x1B)
	fmt.Println(str)
}

func Instance() *LoggerService {
	return loggerSvc
}
