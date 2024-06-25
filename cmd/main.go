package main

import (
	"github.com/text3cn/goodle/goodhttp"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/orm"
)

func main() {
	orm.GetDB()
	goodhttp.Run(Router)
}

func Router(engine *httpserver.Engine) {
	// 我们可以在这里面定义每个请求的路由
	engine.Get("/", Index)
}

func Index(ctx *httpserver.Context) {
	ctx.Resp.Text("Index")
}
