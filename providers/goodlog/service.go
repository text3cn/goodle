package goodlog

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/core"
)

var goodlogSvc *GoodlogService

type GoodlogService struct {
	Service
	c     core.Container
	level byte
}

// 日志级别，只记录大于配置级别的日志
const (
	trace = iota // 最低级别，默认。所有日志，完整链路追踪
	debug        // 开发调试信息
	info         // 业务需要收集的有用信息，例如访客 UA、请求耗时等
	warn         // 警告
	err          // 一般运行时错误
	fatal        // 最高级别，重要性最高，记录导致应用 panic 崩溃的严重错误，
	off          // 日志开关，用于关闭日志的记录
)

type Service interface {
	Trace(output ...interface{})
	Tracef(output ...interface{})
	Debug(output ...interface{})
	Debugf(output ...interface{})
	Info(output ...interface{})
	Infof(output ...interface{})
	Warn(output ...interface{})
	Warnf(output ...interface{})
	Error(output ...interface{})
	Errorf(output ...interface{})
	Fatal(output ...interface{})
	Fatalf(output ...interface{})

	Color(color string, output interface{})
	Colorf(color string, output ...interface{})

	//P(output interface{})
}

// trace
func (self *GoodlogService) Trace(out ...interface{}) {
	if self.level > trace {
		return
	}
	self.output(trace, out...)
}

func (self *GoodlogService) Tracef(out ...interface{}) {
	if self.level > trace {
		return
	}
	self.output(trace, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// debug
func (self *GoodlogService) Debug(out ...interface{}) {
	if self.level > debug {
		return
	}
	self.output(debug, out...)
}

func (self *GoodlogService) Debugf(out ...interface{}) {
	if self.level > debug {
		return
	}
	self.output(debug, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// info
func (self *GoodlogService) Info(out ...interface{}) {
	if self.level > info {
		return
	}
	self.output(info, out...)
}

func (self *GoodlogService) Infof(out ...interface{}) {
	if self.level > info {
		return
	}
	self.output(info, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// warn
func (self *GoodlogService) Warn(out ...interface{}) {
	if self.level > warn {
		return
	}
	self.output(warn, out...)
}

func (self *GoodlogService) Warnf(out ...interface{}) {
	if self.level > warn {
		return
	}
	self.output(warn, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// err
func (self *GoodlogService) Error(out ...interface{}) {
	if self.level > err {
		return
	}
	self.output(err, out...)
}

func (self *GoodlogService) Errorf(out ...interface{}) {
	if self.level > err {
		return
	}
	self.output(err, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// fatal
func (self *GoodlogService) Fatal(out ...interface{}) {
	if self.level > fatal {
		return
	}
	self.output(fatal, out...)
}

func (self *GoodlogService) Fatalf(out ...interface{}) {
	if self.level > fatal {
		return
	}
	self.output(fatal, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}

// color
func (s *GoodlogService) Color(color string, output interface{}) {
	goodlogSvc.P(color, output)
}

func (s *GoodlogService) Colorf(color string, out ...interface{}) {
	goodlogSvc.P(color, fmt.Sprintf(cast.ToString(out[0]), out[1:]...))
}
