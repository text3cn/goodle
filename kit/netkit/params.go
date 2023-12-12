package netkit

import (
	"net/http"
	"strconv"
	"t3im/t3go/logger"
)

type dataConverter struct {
	value string
}

func (this *dataConverter) String(defaultValue ...string) string {
	if this.value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return this.value
}

func (this *dataConverter) Int(defaultValue ...int) int {
	if this.value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	data, _ := strconv.Atoi(this.value)
	return data
}

// 获取请求参数
func Param(r *http.Request, param string) *dataConverter {
	var value string
	if r.Method == "GET" {
		data, ok := r.URL.Query()[param]
		if !ok {
			logger.Cred("获取 GET 参数错误")
		}
		value = data[0]
	}
	return &dataConverter{value: value}
}
