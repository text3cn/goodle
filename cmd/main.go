package main

import (
	"fmt"
	"github.com/text3cn/goodle/kit/gokit"
	"time"
)

func main() {
	var restrictor = gokit.NewTokenBucket(time.Millisecond*10, 100)
	time.Sleep(time.Second * 1) // 填充令牌是异步的，睡眠一会便能获取到令牌
	token := restrictor.TakeToken()
	if token {
		fmt.Println("获得令牌")
	} else {
		fmt.Println("没有令牌，请稍后再试")
	}
}

//
//// 启动 http 服务
//func http() {
//
//	//go goodle.FileServer(8001, "./web")
//
//	goodle.Init().RunHttp(func(engine *httpserver.Engine) {
//		engine.Get("/", Index, middleware.Recovery(map[string]interface{}{}))
//	}, ":3333")
//
//	goodlog.Trace("Trace 级别日志")
//	goodlog.Debug("Debug 级别日志")
//	goodlog.Info("Info 级别日志")
//	goodlog.Warn("Warn 级别日志")
//	goodlog.Error("Error 级别日志")
//	goodlog.Fatal("Fatal 级别日志")
//
//	goodlog.Redf("xxx %d xxx", 100)
//
//	goodle.Init().RunHttp(func(engine *httpserver.Engine) {
//
//		engine.Get("/", Index, middleware.Recovery(map[string]interface{}{}))
//
//	}, ":3333")
//}
//
//func Index(ctx *httpserver.Context) {
//	goodlog.Trace("333")
//	ctx.Resp.Json(map[string]any{
//		"code":    0,
//		"message": "goodle demo",
//	})
//}
//
//func Foo(ctx *httpserver.Context) {
//
//	// 日志属于内置服务，不需要在这实例化，直接框架启动时实例化扔进 context 直接用
//	// IsBind 检查下，不然用户业务服务会覆盖内置服务
//	//time.Sleep(1 * time.Second)
//	//arr := []string{}
//	//println(arr[2])
//	//ctx.NewSingleProvider(cache.Name).(cache.Service).LocalCache("缓存设置")
//	ret := ctx.Req.GetString("a")
//	ctx.Resp.Json(ret)
//	//fmt.Println(c.GetStringSlice("a", []string{}))
//}
