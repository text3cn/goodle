// 分组路由注册
package httpserver

import (
	"strings"
)

// IGroup 路由分组接口
type IGroup interface {
	Get(string, RequestHandler, ...MiddlewareHandler)
	Post(string, RequestHandler, ...MiddlewareHandler)
	Put(string, RequestHandler, ...MiddlewareHandler)
	Delete(string, RequestHandler, ...MiddlewareHandler)
	UseMiddleware(...MiddlewareHandler) IGroup
}

// 实现了 IGroup，按前缀分组
type Prefix struct {
	httpCore *Engine
	prefix   string // 这个group的通用前缀
}

// 初始化前缀分组
func NewPrefix(core *Engine, prefix string) *Prefix {
	return &Prefix{
		httpCore: core,
		prefix:   prefix,
	}
}

func (p *Prefix) Get(uri string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("GET", p.prefix, uri, routeTypeGoGroup, handler, middlewares...)
}

func (p *Prefix) Post(uri string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("POST", p.prefix, uri, routeTypeGoGroup, handler, middlewares...)
}

func (p *Prefix) Put(uri string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("PUT", p.prefix, uri, routeTypeGoGroup, handler, middlewares...)
}

func (p *Prefix) Delete(uri string, handler RequestHandler, middlewares ...MiddlewareHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("DELETE", p.prefix, uri, routeTypeGoGroup, handler, middlewares...)
}

func (p *Prefix) UseMiddleware(middlewares ...MiddlewareHandler) IGroup {
	p.httpCore.groupMiddlewares[p.prefix] = middlewares
	return p
}

// 实现 Group 方法
func (hc *Engine) Prefix(prefix string) IGroup {
	return NewPrefix(hc, prefix)
}
