package httpserver

import (
	"errors"
	"github.com/text3cn/goodle/core"
	"net/http"
	"strings"
)

const (
	routeTypeGoStatic = iota + 1
	routeTypeGoGroup
)

type MiddlewareHandler func(c *Context) error
type RequestHandler func(c *Context) // API / 控制器函数

// 框架核心结构体
type Engine struct {
	router            map[string]map[string]t3WebRoute
	globalMiddlewares []MiddlewareHandler
	groupMiddlewares  map[string][]MiddlewareHandler
	requestHandler    RequestHandler
	container         core.Container
	cross             bool
}

type t3WebRoute struct {
	middlewares    []MiddlewareHandler
	requestHandler RequestHandler
	routeType      int8 // 路由类型：1.golang 静态路由 2.golang 分组路由
	prefix         string
}

// 初始化框架核心结构
func (self *Engine) NewHttpEngine(serviceCenter core.Container) *Engine {
	// 路由 map 的第一维存请求方式，二维存控制器
	router := map[string]map[string]t3WebRoute{}
	router["GET"] = map[string]t3WebRoute{}
	router["POST"] = map[string]t3WebRoute{}
	router["PUT"] = map[string]t3WebRoute{}
	router["DELETE"] = map[string]t3WebRoute{}
	engine := &Engine{
		router:           router,
		groupMiddlewares: map[string][]MiddlewareHandler{}, // 分组路由(批量前缀)路由上挂的中间件
		container:        serviceCenter,
	}
	return engine
}

// 注册全局中间件
func (self *Engine) UseMiddleware(handlers ...MiddlewareHandler) {
	self.globalMiddlewares = handlers
}

// 跨域
func (self *Engine) UseCross() {
	self.cross = true
}

func (self t3WebRoute) IsEmpty() bool {
	return self.routeType == 0
}

// 添加路由到 map
// prefix 主要用于请求进来时匹配分组中间件
func (self *Engine) AddRoute(methoad, prefix, uri string, routeType int8, handler RequestHandler, middlewares ...MiddlewareHandler) {
	uri = strings.ToLower(uri)
	if self.router[methoad][uri].routeType != 0 {
		err := errors.New("route exist: " + uri)
		panic(err)
	}
	self.router[methoad][uri] = t3WebRoute{
		middlewares:    middlewares,
		requestHandler: handler,
		routeType:      routeType,
		prefix:         prefix,
	}
}

func Cross(response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Headers", "*")
}

// 框架核心结构实现 Handler 接口
func (self *Engine) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//self.container.NewSingle(goodlog.Name).(goodlog.Service).Trace("进入 core.serveHTTP")
	//goodlog.Trace("Enter core.serveHTTP")

	if self.cross {
		Cross(response)
	}
	if request.URL.Path == "/favicon.ico" {
		return
	}
	if request.Method == "OPTIONS" {
		response.WriteHeader(200)
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
	// 注入中间件、控制器给 context
	middlewareChain := self.globalMiddlewares
	if route.prefix != "" {
		middlewareChain = append(middlewareChain, self.groupMiddlewares[route.prefix]...)
	}
	for _, middleware := range route.middlewares {
		middlewareChain = append(middlewareChain, middleware)
	}
	ctx.SetMiddwares(middlewareChain)
	// 执行中间件、控制器
	if err := ctx.Next(); err != nil {
		ctx.Resp.SetStatus(500).Text(err.Error())
		return
	}

	// 执行控制器函数
	route.requestHandler(ctx)
}

// 匹配路由，如果没有匹配到，返回 nil
func (this *Engine) FindRouteHandler(request *http.Request) t3WebRoute {
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
