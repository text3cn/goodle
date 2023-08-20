package middleware

import (
	"fmt"
	"github.com/text3cn/goodle/providers/httpserver"
	"log"
	"time"
)

func Cost() httpserver.MiddlewareHandler {
	// 使用函数回调
	return func(c *httpserver.Context) error {
		fmt.Println("use cost middleware")
		// 记录开始时间
		start := time.Now()

		// 使用next执行具体的业务逻辑
		c.Next()

		// 记录结束时间
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.Request().RequestURI, cost.Seconds())

		return nil
	}
}
