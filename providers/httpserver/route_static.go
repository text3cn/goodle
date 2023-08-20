// 静态路由注册
package httpserver

func (this *GoodleEngine) Get(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	this.AddRoute("GET", "", url, routeTypeGoStatic, handler, middlewares...)
}

func (this *GoodleEngine) Post(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	this.AddRoute("POST", "", url, routeTypeGoStatic, handler, middlewares...)
}

func (this *GoodleEngine) Put(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	this.AddRoute("PUT", "", url, routeTypeGoStatic, handler, middlewares...)
}

func (this *GoodleEngine) Delete(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	this.AddRoute("DELETE", "", url, routeTypeGoStatic, handler, middlewares...)
}
