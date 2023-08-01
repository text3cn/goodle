package engine

import (
	"github.com/spf13/cobra"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/logger"
	"github.com/text3cn/goodle/types"
	"log"
	"net/http"
)

// 挂载框架内置命令
func AddKernelCommands(command *Command, router types.HttpEngine) {
	rootCmd := command.rootCmd
	// 后台运行 http 服务
	httpServer := cobra.Command{
		Use:     "start",
		Short:   "Run as a daemon",
		Example: "./main start",
		Run: func(cmd *cobra.Command, args []string) {
			// cmd.Help()
			// 启动子进程
			fork(command, router)
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
			startHttpServer(command, router)
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
func startHttpServer(cmd *Command, router types.HttpEngine) {

	//cntxt := &daemon.Context{
	//	// 设置pid文件
	//	PidFileName: serverPidFile,
	//	PidFilePerm: 0664,
	//	// 设置日志文件
	//	LogFileName: serverLogFile,
	//	LogFilePerm: 0640,
	//	// 设置工作路径
	//	WorkDir: currentFolder,
	//	// 设置所有设置文件的mask，默认为750
	//	Umask: 027,
	//	// 子进程的参数，按照这个参数设置，子进程的命令为 ./main start
	//	Args: []string{"", "start"},
	//}

	container := cmd.container
	// 往 engine 绑定服务是把服务绑定到 http 服务的服务中心
	// 在 http 服务中又另外 new 了一个服务中心，也就是说与框架的服务中心是隔离的
	engine := container.NewSingle(httpserver.Name).(*httpserver.HttpServerService).GoodleEngine.WebServer(container)
	router(engine) // 把路由保存到 map
	cfgsvc := container.NewSingle(config.Name).(config.Service)
	addr := cfgsvc.GetHttpAddr()
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
