// 静态路由注册
package httpserver

func (this *GoodleEngine) Get(url string, handlers ...RequestHandler) {
	this.AddRoute("GET", "", url, routeTypeGoStatic, handlers...)
}

func (this *GoodleEngine) Post(url string, handlers ...RequestHandler) {
	this.AddRoute("POST", "", url, routeTypeGoStatic, handlers...)
}

func (this *GoodleEngine) Put(url string, handlers ...RequestHandler) {
	this.AddRoute("PUT", "", url, routeTypeGoStatic, handlers...)
}

func (this *GoodleEngine) Delete(url string, handlers ...RequestHandler) {
	this.AddRoute("DELETE", "", url, routeTypeGoStatic, handlers...)
}
