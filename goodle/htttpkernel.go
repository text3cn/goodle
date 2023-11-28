package goodle

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/cobra" // https://github.com/spf13/cobra
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/goodlog"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/types"
	"net/http"
)

type HttpEngine func(engine *httpserver.Engine)

type Command struct {
	container core.Container
	rootCmd   *cobra.Command
	config    config.Service
}

func (*Goodle) RunHttp(router HttpEngine, addr ...string) {
	ADDR := ""
	if len(addr) > 0 {
		ADDR = addr[0]
	}
	// 全局容器为框架必要服务，http 容器为用户可选开启 bind 哪些服务
	// 目前服务不多，暂不支持用户自定义服务，所以使用全局服务中心
	c := core.GlobalCore()

	startHttpServer(c, ADDR, router)

	// 容器化运行不需要命令行，将命令行迁移到 good 工具中，通过 good 启动
	//var cobraRoot = &cobra.Command{
	//	// 定义根命令的关键字
	//	Use: "./main",
	//	// 简短介绍
	//	Short: "Goodle Framework",
	//	// 根命令的执行函数
	//	RunE: func(cmd *cobra.Command, args []string) error {
	//		cmd.InitDefaultHelpFlag()
	//		return cmd.Help()
	//	},
	//	// 不需要出现 cobra 默认的 completion 子命令
	//	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	//}
	//var cmd = &Command{
	//	container: c,
	//	rootCmd:   cobraRoot,
	//	config:    config.Instance(),
	//}
	//// 绑定框架内置的命令
	//AddKernelCommands(c, cmd, ADDR, router)

	// 绑定业务的命令
	// AddAppCommand(rootCmd)

	//isDaemon := cmd.config.IsDaemon()
	//if !isDaemon {
	//	// 直接前台挂起运行
	//	startHttpServer(c, ADDR, router)
	//} else {
	//	// 命令行运行，执行 RootCommand
	//	cmd.rootCmd.Execute()
	//}
}

// 挂载框架内置命令
func AddKernelCommands(c *core.ServicesContainer, command *Command, addr string, router HttpEngine) {
	rootCmd := command.rootCmd
	// 后台运行 http 服务
	httpServer := cobra.Command{
		Use:     "start",
		Short:   "Run as a daemon",
		Example: "./main start",
		Run: func(cmd *cobra.Command, args []string) {
			// cmd.Help()
			// 启动子进程
			fork(c, command, addr, router)
			// 绘制主进程控制面板
			drawControl()
		},
	}
	// 子命令，前台运行 http 服务
	httpServer.AddCommand(&cobra.Command{
		Use:     "foreground",
		Short:   "Run in foreground",
		Aliases: []string{"f"},
		Example: "./main start f",
		Run: func(cmd *cobra.Command, args []string) {
			startHttpServer(c, addr, router)
		},
	})
	rootCmd.AddCommand(&httpServer)
	// 停止后台 http 服务
	rootCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stop the daemon",
		Run: func(cmd *cobra.Command, args []string) {
			stop(command)
		},
	})

}

// 启动 http 服务，初始化注册所有内置服务
func startHttpServer(c *core.ServicesContainer, addr string, router HttpEngine) {
	engine := c.NewSingle(httpserver.Name).(*httpserver.HttpServerService).Engine.NewHttpEngine(c)
	router(engine) // 把路由保存到 map
	cfgsvc := c.NewSingle(core.Config).(config.Service)
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
			"/swagger-ui/index.html"
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
