// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package logger

import (
	"github.com/text3cn/goodle/container"
)

const Name = "logger"

type LoggerServiceProvider struct {
	container.ServiceProvider // 显示的写上实现了哪个接口主要是为了代码可读性以及 IDE 友好
}

var loggerSvc *LoggerService

func (self *LoggerServiceProvider) Name() string {
	return Name
}

// 往服务中心注册自己前的操作
func (sp *LoggerServiceProvider) BeforeInit(c container.Container) error {
	return nil
}

// 服务容器通过调用此方法，约定通过 params 第一个参数将服务容器注入给服务提供者
func newLogger(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开，将配置注入到日志类，例如日志开关等
	c := params[0].(container.Container)
	if loggerSvc != nil {
		return loggerSvc, nil
	}
	loggerSvc = &LoggerService{c: c}
	return loggerSvc, nil
}

// 将创建日志服务实例的函数通过回调函数的方式传递给服务中心，
// 这样服务中心就不需要 import 日志服务就持有了日志服务的实例
func (sp *LoggerServiceProvider) RegisterProviderInstance(c container.Container) container.NewInstanceFunc {
	// 初始化实例的方法
	return newLogger
}

// 日志服务不需要延迟初始化，启动程序就需要打印日志了
func (*LoggerServiceProvider) InitOnBoot() bool {
	return true
}

// 实例化的参数
func (sp *LoggerServiceProvider) Params(c container.Container) []interface{} {
	return []interface{}{c}
}

func (*LoggerServiceProvider) AfterInit(instance any) error {
	return nil
}
