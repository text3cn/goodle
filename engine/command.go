package engine

import (
	"github.com/spf13/cobra" // https://github.com/spf13/cobra
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/cache"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/orm"
	"github.com/text3cn/goodle/types"
)

type Command struct {
	container container.Container
	rootCmd   *cobra.Command
	config    config.Service
}

// 初始化服务容器，绑定根 Command 运行
func Run(router types.HttpEngine, beforStrt ...types.BeforStartCallback) {
	c := container.New()
	initServices(c)
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
		config:    c.NewSingle(config.Name).(config.Service),
	}
	// 绑定框架内置的命令
	AddKernelCommands(cmd, router)

	// 绑定业务的命令
	// AddAppCommand(rootCmd)
	if len(beforStrt) > 0 {
		beforStrt[0](c)
	}

	isDevelop := cmd.config.IsDevelop()
	if isDevelop {
		// 直接前台挂起运行
		startHttpServer(cmd, router)
	} else {
		// 命令行运行，执行 RootCommand
		cmd.rootCmd.Execute()
	}

}

func initServices(c container.Container) {
	c.Bind(&httpserver.HttpServerProvider{})
	c.Bind(&config.ConfigProvider{})
	c.Bind(&cache.CacheServiceProvider{})
	c.Bind(&orm.OrmProvider{})
	ormService := c.NewSingle(orm.Name).(orm.Service)
	ormService.Init()
}
