package httpserver

import (
	"context"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/kit/castkit"
	"github.com/text3cn/goodle/providers/cache"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/redis"
	"net/http"
	"sync"
	"time"
)

// 实现标准库的 Context
// 基本现在所有第三方库函数都会根据官方的建议将第一个参数设置为标准 Context 接口，
// 所以定制一个自己的 context 很有用，这里将 request 和 response 封装到 context，
// 这样就可以在整条请求链路中随时处理输入输出
type Context struct {
	request        *http.Request
	context        context.Context
	middwares      []MiddlewareHandler // 中间件
	middwaresIndex int                 // 用数组加索引偏移来实现中间件到控制器的调用链
	// 边界场景处理：
	// 异常、超时事件触发时，需要往 responseWriter 中写入信息给客户端，
	// 这时候如果有其他 Goroutine 也在操作 responseWriter 可能会出现 responseWriter 中的信息重复写入，
	// 并且写入的顺序也可能是错误乱的，分两步解决：
	// 1. 写保护，在写 response 的时候加锁，保证顺序正确
	writerMux *sync.Mutex
	// 2. 添加标记，当发生 timeout 时设置标记位为 true，在 Context 提供的写 response 函数中，
	//    先读取标记位，如果为 true，表示已经给客户端返回过了，就不要再写 response 了。
	hasTimeout bool
	// 服务中心
	container core.Container
	values    map[string]interface{}

	// 配置服务
	Req    IRequest
	Resp   RespStruct
	Config config.Service
	Cache  cache.Service
	Redis  redis.Service
}

func NewContext(r *http.Request, w http.ResponseWriter, holder core.Container) *Context {
	req := &ReqStruct{request: r}
	ctx := &Context{
		context:        r.Context(),
		writerMux:      &sync.Mutex{},
		middwaresIndex: -1,
		container:      holder,
		values:         map[string]interface{}{},

		Req:    req,
		Resp:   RespStruct{request: req, responseWriter: w},
		Config: holder.NewSingle(core.Config).(config.Service),
		Cache:  holder.NewSingle(cache.Name).(cache.Service),
		Redis:  holder.NewSingle(redis.Name).(redis.Service),
	}
	return ctx
}

// 对用户暴露服务者中心
func (ctx *Context) Holder() core.Container {
	return ctx.container
}

// 对外暴露锁
func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

// 请求时中间件
func (ctx *Context) SetMiddwares(handlers []MiddlewareHandler) {
	ctx.middwares = handlers
}

// 按顺序执行中间件
func (ctx *Context) Next() error {
	// 执行中间件
	ctx.middwaresIndex++
	if ctx.middwaresIndex < len(ctx.middwares) {
		if err := ctx.middwares[ctx.middwaresIndex](ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.Resp.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.context
}

// #region implement context.Context
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.context.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.context.Done()
}

func (ctx *Context) Err() error {
	return ctx.context.Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.context.Value(key)
}

// 将服务注册到服务中心
func (ctx *Context) NewSingleProvider(name string) interface{} {
	return ctx.container.NewSingle(name)
}

func (ctx *Context) NewInstanceProvider(name string, params ...interface{}) interface{} {
	return ctx.container.NewInstance(name, params)
}

// 往 context 上设置值/获取值
func (ctx *Context) SetVal(key string, value interface{}) {
	ctx.values[key] = value
}
func (ctx *Context) GetVal(key string) *castkit.GoodleVal {
	return &castkit.GoodleVal{ctx.values[key]}
}
