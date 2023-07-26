// 分组路由注册
package httpserver

import "strings"

// IGroup 路由分组接口
type IGroup interface {
	Get(string, ...RequestHandler)
	Post(string, ...RequestHandler)
	Put(string, ...RequestHandler)
	Delete(string, ...RequestHandler)
	UseMiddleware(...RequestHandler) IGroup
}

// 实现了 IGroup，按前缀分组
type Prefix struct {
	httpCore *T3WebEngine
	prefix   string // 这个group的通用前缀
}

// 初始化前缀分组
func NewPrefix(core *T3WebEngine, prefix string) *Prefix {
	return &Prefix{
		httpCore: core,
		prefix:   prefix,
	}
}

func (p *Prefix) Get(uri string, handlers ...RequestHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("GET", p.prefix, uri, routeTypeGoGroup, handlers...)
}

func (p *Prefix) Post(uri string, handlers ...RequestHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("POST", p.prefix, uri, routeTypeGoGroup, handlers...)
}

func (p *Prefix) Put(uri string, handlers ...RequestHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("PUT", p.prefix, uri, routeTypeGoGroup, handlers...)
}

func (p *Prefix) Delete(uri string, handlers ...RequestHandler) {
	uri = strings.ToLower(p.prefix + uri)
	p.httpCore.AddRoute("DELETE", p.prefix, uri, routeTypeGoGroup, handlers...)
}

func (p *Prefix) UseMiddleware(handlers ...RequestHandler) IGroup {
	p.httpCore.groupMiddlewares[p.prefix] = handlers
	return p
}

// 实现 Group 方法
func (hc *T3WebEngine) Prefix(prefix string) IGroup {
	return NewPrefix(hc, prefix)
}
