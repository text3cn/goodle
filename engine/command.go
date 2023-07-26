package engine

import (
	"github.com/spf13/cobra" // https://github.com/spf13/cobra
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/httpserver"
)

type Command struct {
	container container.Container
	cobra     *cobra.Command
}

// 初始化根 Command 并运行
func Run(container container.Container, router func(engine *httpserver.GoodleEngine)) {
	var cobraRoot = &cobra.Command{
		// 定义根命令的关键字
		Use: "./main",
		// 简短介绍
		Short: "Goodle Command",
		// 根命令的详细介绍
		Long: "Goodle Framework For Command",
		// 根命令的执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现 cobra 默认的 completion 子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	var cmd = &Command{
		container: container,
		cobra:     cobraRoot,
	}
	// 绑定框架内置的命令
	AddKernelCommands(cmd, router)
	// 绑定业务的命令
	// AddAppCommand(rootCmd)

	httpServerDeamon(cmd, router) // 前台运行

	// 执行 RootCommand
	cmd.cobra.Execute()
}

func (self *Command) GetContainer() container.Container {
	return self.container
}
