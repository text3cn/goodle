package middleware

import (
	"context"
	"fmt"
	"github.com/text3cn/goodle/providers/httpserver"
	"log"
	"time"
)

// http 请求超时控制的中间件
func Timeout(timeout time.Duration) httpserver.MiddlewareHandler {
	// 使用回调函数，返回一个匿名函数保存到 handlers，由 context.Next() 进行调用
	return func(ctx *httpserver.Context) error {
		fmt.Println("use timeout middleware")
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 执行业务逻辑前预操作：创建超时 context
		durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), timeout)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 继续往下执行中间件或业务逻辑
			ctx.Next()
			finish <- struct{}{}
		}()
		// 信号监听
		select {
		case err := <-panicChan:
			//ctx.Json(500, err)
			panic(err)
		case <-durationCtx.Done():
			// 业务处理超时
			log.Println("超时")
			ctx.WriterMux().Lock()
			defer ctx.WriterMux().Unlock()
			ctx.Resp.SetStatus(500).Json("time out")
			ctx.SetHasTimeout()
		case <-finish:
			// 业务正常处理完毕
			fmt.Println("finish")
		}
		return nil
	}
}
