// 日志服务提供者对外提供的能力
package logger

import (
	"github.com/text3cn/goodle/container"
)

// 声明对外服务的接口
type Service interface {
	Trace(output string) // 跟踪程序运行轨迹
	Debug(output string)
	Fatal(output string) // panic 级错误
	// INFO
	// WARN
	Error(output string)

	Fatalf(output string) string
}

// 定义日志类
type LoggerService struct {
	Service // 实现接口
	c       container.Container
}

// 实现接口
func (s *LoggerService) Trace(out string) {
	println("[TRACE]", out)
}

func (s *LoggerService) Debug(out string) {
	println("[DEBUG]", out)
}

func (s *LoggerService) Fatal(out string) {
	println("[FATAL]", out)
}

func (s *LoggerService) Error(out string) {
	println("[Error]", out)
}

func (s *LoggerService) Fatalf(out string) string {
	return "[Fatalf]" + out
}
