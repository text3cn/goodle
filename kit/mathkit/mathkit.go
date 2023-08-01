package mathkit

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// 取随机数
// 范围: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// 对浮点数进行四舍五入
// args[0] 保留几位小数
func Round(value float64, args ...int) float64 {
	if len(args) > 0 {
		value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
		return value
	}
	return math.Floor(value + 0.5)
}
