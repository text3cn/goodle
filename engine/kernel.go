package engine

import (
	"github.com/spf13/cobra"
	"github.com/text3cn/goodle/providers/cache"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/logger"
	"log"
	"net/http"
)

// 挂载框架内置命令
func AddKernelCommands(command *Command, router func(engine *httpserver.GoodleEngine)) {
	httpServer := &cobra.Command{
		Use:   "start",
		Short: "Start as a daemon",
		Run: func(cmd *cobra.Command, args []string) {
			httpServerDeamon(command, router)
		},
	}
	command.cobra.AddCommand(httpServer)
}

func httpServerDeamon(command *Command, router func(engine *httpserver.GoodleEngine)) {
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
	startHttpServer(command, router)
}

func startHttpServer(command *Command, router func(engine *httpserver.GoodleEngine)) {
	container := command.GetContainer()
	container.Bind(&httpserver.HttpServerProvider{})
	engine := container.NewSingle(httpserver.Name).(*httpserver.HttpServerService).GoodleEngine.WebServer()
	engine.ServiceProvider(logger.Name, &logger.LoggerServiceProvider{})
	engine.ServiceProvider(cache.Name, &cache.CacheServiceProvider{})
	router(engine) // 把路由保存到 map
	server := &http.Server{
		// 自定义的请求核心处理函数
		Handler: engine,
		// 请求监听地址
		Addr: ":9000",
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Println("[Start http fail]", err)
	}
}
