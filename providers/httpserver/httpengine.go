package httpserver

import (
	"errors"
	servicescenter2 "github.com/text3cn/t3web/container"
	"github.com/text3cn/t3web/providers/logger"
	"net/http"
	"strings"
)

const (
	routeTypeGoStatic = iota + 1
	routeTypeGoGroup
)

type RequestHandler func(c *Context) error

// 框架核心结构体
type T3WebEngine struct {
	router            map[string]map[string]t3WebRoute
	globalMiddlewares []RequestHandler
	groupMiddlewares  map[string][]RequestHandler
	container         servicescenter2.Container
}

type t3WebRoute struct {
	handlers  []RequestHandler
	routeType int8 // 路由类型：1.golang 静态路由 2.golang 分组路由
	prefix    string
}

// 初始化框架核心结构
func (self *T3WebEngine) T3Web() *T3WebEngine {
	// 路由 map 的第一维存请求方式，二维存控制器
	router := map[string]map[string]t3WebRoute{}
	router["GET"] = map[string]t3WebRoute{}
	router["POST"] = map[string]t3WebRoute{}
	router["PUT"] = map[string]t3WebRoute{}
	router["DELETE"] = map[string]t3WebRoute{}
	engine := &T3WebEngine{
		router:           router,
		groupMiddlewares: map[string][]RequestHandler{},
		container:        servicescenter2.New(),
	}
	return engine
}

// 注册全局中间件
func (this *T3WebEngine) UseMiddleware(handlers ...RequestHandler) {
	this.globalMiddlewares = handlers
}

func (this t3WebRoute) IsEmpty() bool {
	return this.routeType == 0
}

// 添加路由到 map
// prefix 主要用于请求进来时匹配分组中间件
func (this *T3WebEngine) AddRoute(methoad, prefix, uri string, routeType int8, handlers ...RequestHandler) {
	uri = strings.ToLower(uri)
	if this.router[methoad][uri].routeType != 0 {
		err := errors.New("route exist: " + uri)
		panic(err)
	}
	handlerChain := append(handlers[1:], handlers[0]) // 把响应函数移到中间件后面
	this.router[methoad][uri] = t3WebRoute{
		handlers:  handlerChain,
		routeType: routeType,
		prefix:    prefix,
	}
}

// 框架核心结构实现 Handler 接口
func (self *T3WebEngine) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	self.container.NewSingle(logger.Name).(logger.Service).Trace("进入 core.serveHTTP")
	if request.URL.Path == "/favicon.ico" {
		return
	}
	// 初始化自定义 context
	ctx := NewContext(request, response, self.container)
	// 寻找路由，handlers 包含中间件 + 控制器
	route := self.FindRouteHandler(request)
	if route.IsEmpty() {
		ctx.Resp.SetStatus(404).Text("404 not found")
		return
	}

	// 添加响应函数的调用链到 context，由 context 发起调用
	handlersChain := self.globalMiddlewares
	if route.prefix != "" {
		handlersChain = append(handlersChain, self.groupMiddlewares[route.prefix]...)
	}
	handlersChain = append(handlersChain, route.handlers...)
	ctx.SetHandlers(handlersChain)

	// 执行路由绑定的函数，如果返回 err 代表存在内部错误，返回 500 状态码
	if err := ctx.Next(); err != nil {
		ctx.Resp.SetStatus(500).Text(err.Error())
		return
	}
}

// 匹配路由，如果没有匹配到，返回 nil
func (this *T3WebEngine) FindRouteHandler(request *http.Request) t3WebRoute {
	// 转换大小写，确保大小写不敏感
	method := strings.ToUpper(request.Method)
	uri := strings.ToLower(request.URL.Path)
	// 查找第一层map
	if methodHandlers, ok := this.router[method]; ok {
		if handler, ok := methodHandlers[uri]; ok {
			return handler
		}
	}
	return t3WebRoute{}
}

// 服务注册其实跟 http 路由一样的原理
func (self *T3WebEngine) ServiceProvider(name string, provider servicescenter2.ServiceProvider) {
	self.container.Bind(provider)
}

// 静态资源服务
//fs := http.FileServer(http.Dir("/home/bob/static"))
//http.Handle("/static/", http.StripPrefix("/static", fs))
