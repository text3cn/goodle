package engine

import (
	"github.com/spf13/cobra" // https://github.com/spf13/cobra
	"github.com/text3cn/goodle/config"
	config2 "github.com/text3cn/goodle/config"
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/logger"
	"log"
	"net/http"
)

type HttpEngine func(engine *httpserver.GoodleEngine)

type Command struct {
	container container.Container
	rootCmd   *cobra.Command
	config    config.Service
}

var beforStartInvoke config2.BeforStartCallback

type runHttp struct {
	beforStart config2.BeforStartCallback
}

func BeforStartHttp(beforStart config2.BeforStartCallback) *runHttp {
	beforStartInvoke = beforStart
	return &runHttp{beforStart: beforStart}
}

func (*runHttp) RunHttp(router HttpEngine, addr ...string) {
	RunHttp(router, addr...)
}

// 初始化服务容器，绑定根 Command 运行
func RunHttp(router HttpEngine, addr ...string) {
	ADDR := ""
	if len(addr) > 0 {
		ADDR = addr[0]
	}
	c := container.New()
	bindServices(c)
	var cobraRoot = &cobra.Command{
		// 定义根命令的关键字
		Use: "./main",
		// 简短介绍
		Short: "Goodle Framework",
		// 根命令的执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现 cobra 默认的 completion 子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	var cmd = &Command{
		container: c,
		rootCmd:   cobraRoot,
		config:    config.Instance(),
	}
	// 绑定框架内置的命令
	AddKernelCommands(c, cmd, ADDR, router)

	// 绑定业务的命令
	// AddAppCommand(rootCmd)
	if beforStartInvoke != nil {
		beforStartInvoke(c)
	}
	isDaemon := cmd.config.IsDaemon()
	if !isDaemon {
		// 直接前台挂起运行
		startHttpServer(c, ADDR, router)
	} else {
		// 命令行运行，执行 RootCommand
		cmd.rootCmd.Execute()
	}
}

// 挂载框架内置命令
func AddKernelCommands(c *container.ServicesContainer, command *Command, addr string, router HttpEngine) {
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
func startHttpServer(c *container.ServicesContainer, addr string, router HttpEngine) {
	engine := c.NewSingle(httpserver.Name).(*httpserver.HttpServerService).GoodleEngine.WebServer(c)
	router(engine) // 把路由保存到 map
	cfgsvc := config.Instance()
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
	logger.Instance().Trace("Server Listen On " + addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println("[Start http fail]", err)
	}
}
