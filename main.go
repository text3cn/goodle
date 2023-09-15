package main

import (
	"github.com/text3cn/goodle/engine"
	"github.com/text3cn/goodle/providers/cache"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/httpserver/middleware"
	"time"
)

func main() {

	// 存储桶容量，字节为单位，最小 32MB ，少于 32MB 会当做 32MB 处理
	bucketSize := 32 * 1024 * 1024
	cache.NewFreeCache("bucket1", bucketSize)

	// 启动 http 服务
	engine.RunHttp(func(engine *httpserver.GoodleEngine) {
		engine.Get("/", Index)
	}, ":3333")
}

func Index(ctx *httpserver.Context) {
	ctx.Redis.Pipeline("key1").Redis(300, "slave").Redis(3600*24, "master")

}

func router(core *httpserver.GoodleEngine) {

	// 设置控制器
	// core.Get("foo", framework.TimeoutHandler(FooControllerHandler, time.Second*1))

	// 使用全局中间件，对所有路由生效
	core.UseMiddleware(
		// 开发环境不许要 recover 会把错误栈给吃掉
		//middleware.Recovery(map[string]interface{}{
		//	"msg": "Server Error",
		//}),
		//middleware.Logger(),
		//middleware.Timeout(2*time.Second),
	)

	// 静态路由
	core.Get("/foo", Foo)
	core.Get("/bar", Bar)

	// 批量路由前缀
	prefix := core.Prefix("/controlleR").UseMiddleware(middleware.Cost())
	{
		prefix.Get("/action1", Foo, middleware.Cost())
		prefix.Get("/action2", Action2)
	}
}

func Foo(ctx *httpserver.Context) {

	// 日志属于内置服务，不需要在这实例化，直接框架启动时实例化扔进 context 直接用
	// IsBind 检查下，不然用户业务服务会覆盖内置服务
	//time.Sleep(1 * time.Second)
	//arr := []string{}
	//println(arr[2])
	//ctx.NewSingleProvider(cache.Name).(cache.Service).LocalCache("缓存设置")
	ret, _ := ctx.Req.GetString("a")
	ctx.Resp.Json(ret)
	//fmt.Println(c.GetStringSlice("a", []string{}))
}

func Bar(c *httpserver.Context) {
	c.Resp.Json("bar")

}

func Action2(c *httpserver.Context) {
	time.Sleep(0 * time.Second)
	c.Resp.Json("ok2")
}
