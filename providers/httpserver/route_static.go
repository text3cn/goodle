// 静态路由注册
package httpserver

func (this *T3WebEngine) Get(url string, handlers ...RequestHandler) {
	this.AddRoute("GET", "", url, routeTypeGoStatic, handlers...)
}

func (this *T3WebEngine) Post(url string, handlers ...RequestHandler) {
	this.AddRoute("POST", "", url, routeTypeGoStatic, handlers...)
}

func (this *T3WebEngine) Put(url string, handlers ...RequestHandler) {
	this.AddRoute("PUT", "", url, routeTypeGoStatic, handlers...)
}

func (this *T3WebEngine) Delete(url string, handlers ...RequestHandler) {
	this.AddRoute("DELETE", "", url, routeTypeGoStatic, handlers...)
}
