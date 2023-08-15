package httpserver

import (
	"errors"
	"github.com/text3cn/goodle/container"
	"net/http"
	"strings"
)

const (
	routeTypeGoStatic = iota + 1
	routeTypeGoGroup
)

type RequestHandler func(c *Context) error

// 框架核心结构体
type GoodleEngine struct {
	router            map[string]map[string]t3WebRoute
	globalMiddlewares []RequestHandler
	groupMiddlewares  map[string][]RequestHandler
	container         container.Container
	cross             bool
}

type t3WebRoute struct {
	handlers  []RequestHandler
	routeType int8 // 路由类型：1.golang 静态路由 2.golang 分组路由
	prefix    string
}

// 初始化框架核心结构
func (self *GoodleEngine) WebServer(serviceCenter container.Container) *GoodleEngine {
	// 路由 map 的第一维存请求方式，二维存控制器
	router := map[string]map[string]t3WebRoute{}
	router["GET"] = map[string]t3WebRoute{}
	router["POST"] = map[string]t3WebRoute{}
	router["PUT"] = map[string]t3WebRoute{}
	router["DELETE"] = map[string]t3WebRoute{}
	engine := &GoodleEngine{
		router:           router,
		groupMiddlewares: map[string][]RequestHandler{},
		container:        serviceCenter,
	}
	return engine
}

// 注册全局中间件
func (self *GoodleEngine) UseMiddleware(handlers ...RequestHandler) {
	self.globalMiddlewares = handlers
}

// 跨域
func (self *GoodleEngine) UseCross() {
	self.cross = true
}

func (self t3WebRoute) IsEmpty() bool {
	return self.routeType == 0
}

// 添加路由到 map
// prefix 主要用于请求进来时匹配分组中间件
func (self *GoodleEngine) AddRoute(methoad, prefix, uri string, routeType int8, handlers ...RequestHandler) {
	uri = strings.ToLower(uri)
	if self.router[methoad][uri].routeType != 0 {
		err := errors.New("route exist: " + uri)
		panic(err)
	}
	handlerChain := append(handlers[1:], handlers[0]) // 把响应函数移到中间件后面
	self.router[methoad][uri] = t3WebRoute{
		handlers:  handlerChain,
		routeType: routeType,
		prefix:    prefix,
	}
}

func Cross(response http.ResponseWriter) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Headers", "*")
}

// 框架核心结构实现 Handler 接口
func (self *GoodleEngine) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//self.container.NewSingle(logger.Name).(logger.Service).Trace("进入 core.serveHTTP")
	//logger.Instance().Trace("Enter core.serveHTTP")
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
func (this *GoodleEngine) FindRouteHandler(request *http.Request) t3WebRoute {
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

// 静态资源服务
//fs := http.FileServer(http.Dir("/home/bob/static"))
//http.Handle("/static/", http.StripPrefix("/static", fs))
