package goodhttp

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/goodlog"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/i18n"
	"github.com/text3cn/goodle/types"
	"net/http"
)

type HttpEngine func(engine *httpserver.Engine)

func Run(router HttpEngine, addr ...string) {
	ADDR := ""
	if len(addr) > 0 {
		ADDR = addr[0]
	}
	// 服务容器
	c := core.NewContainer()
	c.Bind(&config.ConfigProvider{})
	c.Bind(&goodlog.GoodlogProvider{})
	c.Bind(&i18n.I18nProvider{})
	c.Bind(&httpserver.HttpServerProvider{})
	startHttpServer(c, ADDR, router)
}

// 启动 http 服务
func startHttpServer(c *core.ServicesContainer, addr string, router HttpEngine) {
	cfgsvc := c.NewSingle(config.Name).(config.Service)
	engine := c.NewSingle(httpserver.Name).(*httpserver.HttpServerService).Engine.NewHttpEngine(c)
	router(engine) // 把路由保存到 map

	// 代码中没有传递端口则去配置文件找
	if addr == "" {
		addr = cfgsvc.GetHttpAddr()
	}
	server := &http.Server{
		// 自定义的请求核心处理函数
		Handler: engine,
		// 请求监听地址
		Addr: addr,
	}
	httpServerOutput(cfgsvc, addr)
	err := server.ListenAndServe()
	if err != nil {
		goodlog.Error("[Start http fail]", err)
	}
}

func httpServerOutput(cfgsvc config.Service, addr string) {
	// web server
	info := fmt.Sprintf("\033[36m%s"+"\033[0m", "WebServer: http://localhost"+addr)
	fmt.Println("")
	fmt.Println(info)
	// swager server
	swagCfg := cfgsvc.GetSwagger()
	if swagCfg != (types.SwaggerConfig{}) {
		str := "SwaggerUI: http://" + swagCfg.SwaggerUiHost + ":" + cast.ToString(swagCfg.SwaggerUiPort) +
			"/swagger-ui/index.html\n"
		info = fmt.Sprintf("\033[36m%s"+"\033[0m", str)
		fmt.Println(info)
	}
}

// dir 相对于可执行文件的当前目录
func FileServer(host string, dir string) {
	// 静态文件服务器
	var staticServer = func(w http.ResponseWriter, req *http.Request) {
		var staticHandler = http.FileServer(http.Dir(dir))
		if req.URL.Path == "/" {
			// 直接访问文件服务器的根目录会递归出所有文件，这里处理成访问根目录时返回自定义的 404 页面
			//req.URL.Path = "/index.html"
		}
		staticHandler.ServeHTTP(w, req)
	}
	// 把跟路径拿来做静态资源服务器
	http.HandleFunc("/", staticServer)
	// 监听端口启动服务
	err := http.ListenAndServe(host, nil)
	if err != nil {
		fmt.Println("http listen failed")
	}
}
