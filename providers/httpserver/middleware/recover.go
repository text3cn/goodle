package middleware

import (
	"errors"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/httpserver"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

// recovery 机制，对协程中的函数异常进行捕获，这个应该作为最外层，即第一个被调用的中间件
func Recovery(diyMsg map[string]interface{}) httpserver.MiddlewareHandler {
	calback := func(c *httpserver.Context) error {
		cfgSvc := c.Holder().NewSingle(core.Config).(config.Service)
		if cfgSvc.IsDebug() {
			return nil
		}
		println("use recovery middleware")
		// 捕获 c.Next() 出现的panic
		isAbort := false
		defer func() {
			if err := recover(); err != nil {
				log.Println("[Recovery] ERR:", err)
				// 底层连接也可能会出现异常，如果持续给已经中断的连接发送请求，会在底层持续显示网络连接错误（broken pipe）
				// 判断是否是底层连接异常，如果是的话，则标记 brokenPipe
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if brokenPipe {
					isAbort = true
				}
				// TODO 打印错误日志，参考 GIN 有细节：https://time.geekbang.org/column/article/422990
				ret := map[string]interface{}{}
				// 用户自定义错误消息
				if diyMsg != nil {
					if orignError, ok := err.(runtime.Error); ok {
						// 追加原始错误消息
						diyMsg["err"] = orignError.Error()
					}
					ret = diyMsg
				} else {
					// 默认错误消息格式
					ret = map[string]interface{}{
						"err": err,
					}
				}
				c.Resp.SetStatus(500).Json(ret)
			}
		}()
		if isAbort {
			return errors.New("broken pipe")
		}
		c.Next()
		return nil

	}
	return calback
}
