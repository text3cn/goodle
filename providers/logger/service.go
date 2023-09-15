// 日志服务提供者对外提供的能力
package logger

import (
	"github.com/text3cn/goodle/container"
)

const (
	OFF   = iota // 最高级别，用于关闭日志记录
	FATAL        // 导致应用崩溃的严重错误
	ERROR        // 其他运行时错误
	WARN         // 警告
	INFO
	DEBUG
	TRACE
)

// 声明对外服务的接口
type Service interface {
	Trace(output ...interface{}) // 跟踪程序运行轨迹
	Debug(output ...interface{})
	Fatal(output ...interface{}) // panic 级错误
	Info(output ...interface{})
	Warn(output ...interface{})
	Error(output ...interface{})
	Fatalf(output ...interface{}) string
	Pink(output ...interface{})
}

// 定义日志类
type LoggerService struct {
	Service // 实现接口
	c       container.Container
}

// 实现接口
func (s *LoggerService) Trace(out ...interface{}) {
	println("[TRACE]", out[0].(string))
}

func (s *LoggerService) Debug(out ...interface{}) {
	println("[DEBUG]", out)
}

func (s *LoggerService) Fatal(out ...interface{}) {
	println("[FATAL]", out)
}

func (s *LoggerService) Info(out ...interface{}) {
	println("[Info]", out)
}

func (s *LoggerService) Warn(out ...interface{}) {
	println("[Warn]", out)
}

func (s *LoggerService) Error(out ...interface{}) {
	println("[Error]", out[0].(string))
}

func (s *LoggerService) Fatalf(out ...interface{}) string {
	return "[Fatalf]"
}

func (s *LoggerService) Pink(output ...interface{}) {
	outputWithColor(35, output...)
}
