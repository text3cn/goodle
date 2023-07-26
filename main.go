package main

import (
	"github.com/text3cn/t3web/container"
	"github.com/text3cn/t3web/engine"
	"github.com/text3cn/t3web/providers/cache"
	"github.com/text3cn/t3web/providers/httpserver"
	"github.com/text3cn/t3web/providers/httpserver/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go engine.Run(container.New(), registerRouter)
	// 优雅关闭，退出进程有以下四种信号：
	// SIGINT  : 前台运行模式下 Windows/Linux 都可以通过 Ctrl+C 键来产生 SIGINT 信号请求中断进程
	// SIGQUIT : 与 SIGINT 类似，前台模式下 Ctrl+\ 通知进程中断，唯一不同是默认会产生 core 文件
	// SIGTERM : 通过 kill pid 命令结束后台进程
	// SIGKILL : 通过 kill -9 pid 命令强制结束后台进程
	// 除了 SIGKILL 信号无法被 Golang 捕获，其余三个信号都是能被阻塞和捕获的。
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) // 订阅指定信号
	<-quit                                                                // 阻塞当前 Goroutine 等待订阅的中断信号
	// 调用Server.Shutdown graceful结束()
	// httpcore.Serve 可以返回回调函数在这里启动服务，然后 Shutdown
	//if err := server.Shutdown(context.Background()); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	// 验证, 在控制器写个 time.Sleep(10 * time.Second) ，然后浏览器访问，然后杀进程看会不会不会等 10 秒睡完再结束
	// 还需要超时处理，如果控制器中睡 1 个小时或者死循环那就无法退出了
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatal("server shutdown error:", err)
	//}
	//select {
	//case <-ctx.Done():
	//	log.Println("timeout of 5 seconds")
	//}
	//log.Println("server exiting")
}

func registerRouter(core *httpserver.T3WebEngine) {
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

func Foo(ctx *httpserver.Context) error {
	// 日志属于内置服务，不需要在这实例化，直接框架启动时实例化扔进 context 直接用
	// IsBind 检查下，不然用户业务服务会覆盖内置服务
	//time.Sleep(1 * time.Second)
	//arr := []string{}
	//println(arr[2])
	ctx.Logger.Trace("wowowoow")
	ctx.NewSingleProvider(cache.Name).(cache.Service).LocalCache("缓存设置")
	ret, _ := ctx.Req.GetString("a")
	ctx.Resp.Json(ret)
	//fmt.Println(c.GetStringSlice("a", []string{}))

	return nil
}

func Bar(c *httpserver.Context) error {
	c.Resp.Json("bar")
	return nil
}

func Action2(c *httpserver.Context) error {
	time.Sleep(0 * time.Second)
	c.Resp.Json("ok2")
	return nil
}
