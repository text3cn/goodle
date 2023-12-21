package main

import (
	"github.com/text3cn/goodle/goodle"
	"github.com/text3cn/goodle/providers/etcd"
)

func main() {
	goodle.Init()
	etcd.Instance().ServiceRegister()
	select {}
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
