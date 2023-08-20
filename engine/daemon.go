package engine

import (
	"github.com/sevlyar/go-daemon"
	"github.com/text3cn/goodle/kit/strkit"
	"github.com/text3cn/goodle/providers/logger"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// daemon 启动成功后父进程 return，子进程运行时脱离控制台，
// 子进程向控制台打印日志时，会被定向到 /dev/null 所以控制台是没有输出的，
// 因此需要将子进程的输出保存到文件中。
func fork(command *Command, addr string, router HttpEngine) {
	processName := strkit.StrReplace("./", "", os.Args[0], 1)
	runtimePath := command.config.GetRuntimePath()
	ctx := &daemon.Context{
		PidFileName: filepath.Join(runtimePath, processName+".pid"),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(runtimePath, processName+".log"),
		LogFilePerm: 0640,
		Umask:       027,
		WorkDir:     "./",                  // 子进程工作目录
		Args:        []string{"", "start"}, // 传递给子进程的参
	}
	// 拷贝上下文创建子进程
	child, err := ctx.Reborn()
	defer ctx.Release()
	if err != nil {
		panic("Failed to create Child process, error:" + err.Error())
	}
	if child != nil {
		// 父进程工作完成，直接在这退出了
		return
	}
	// 子进程启动 http 服务
	go startHttpServer(command, addr, router)
	// 优雅关闭，退出进程有以下四种信号：
	// SIGINT  : 前台运行模式下 Windows/Linux 都可以通过 Ctrl+C 键来产生 SIGINT 信号请求中断进程
	// SIGQUIT : 与 SIGINT 类似，前台模式下 Ctrl+\ 通知进程中断，唯一不同是默认会产生 core 文件
	// SIGTERM : 通过 kill pid 命令结束后台进程
	// SIGKILL : 通过 kill -9 pid 命令强制结束后台进程
	// 除了 SIGKILL 信号无法被 Golang 捕获，其余三个信号都是能被阻塞和捕获的。
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) // 订阅指定信号
	<-quit                                                                // 阻塞当前 Goroutine 等待订阅的中断信号
	// TODO 调用Server.Shutdown graceful结束()
	// httpcore.Serve 可以返回回调函数在这里启动服务，然后 Shutdown
	//if err := server.Shutdown(context.Background()); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	// 验证, 在控制器写个 time.Sleep(10 * time.Second) ，然后浏览器访问，然后杀进程看会不会不会等 10 秒睡完再结束
	// 还需要超时处理，如果控制器中睡 1 个小时或者死循环那就无法退出了
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatal("server shutdown error:", err)
	//}
	//select {
	//case <-ctx.Done():
	//	log.Println("timeout of 5 seconds")
	//}
	//log.Println("server exiting")
}

// 停止子进程
func stop(command *Command) {
	var pid int
	var err error
	var process *os.Process
	processName := strkit.StrReplace("./", "", os.Args[0], 1)
	runtimePath := command.config.GetRuntimePath()
	pidfile := filepath.Join(runtimePath, processName+".pid")
	if pid, err = daemon.ReadPidFile(pidfile); err != nil {
		logger.Instance().Error("pid not found")
		return
	}
	process, err = os.FindProcess(pid)    // 通过 pid 获取子进程
	err = process.Signal(syscall.SIGTERM) // 给子进程发送中断信号
	if err != nil {
		logger.Instance().Error("Stop daemon fail" + err.Error())
	} else {
		logger.Instance().Trace("Stop daemon success")
	}
}
