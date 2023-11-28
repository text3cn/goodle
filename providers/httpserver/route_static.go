// 静态路由注册
package httpserver

func (this *Engine) Get(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	this.AddRoute("GET", "", url, routeTypeGoStatic, handler, middlewares...)
}

func (this *Engine) Post(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {

	this.AddRoute("POST", "", url, routeTypeGoStatic, handler, middlewares...)
}

func (this *Engine) Put(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	this.AddRoute("PUT", "", url, routeTypeGoStatic, handler, middlewares...)
}

func (this *Engine) Delete(url string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	this.AddRoute("DELETE", "", url, routeTypeGoStatic, handler, middlewares...)
}
