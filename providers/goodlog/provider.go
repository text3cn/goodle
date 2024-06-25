// 实现服务中心规定的服务注册要求，遵循注册协议 engine.Container
package goodlog

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"strings"
)

const Name = "goodlog"

type GoodlogProvider struct {
	core.ServiceProvider // 显示的写上实现了哪个接口主要是为了代码可读性以及 IDE 友好
	level                byte
}

func (self *GoodlogProvider) Name() string {
	return Name
}

// 日志服务不需要延迟初始化，启动程序就需要打印日志了
func (*GoodlogProvider) InitOnBoot() bool {
	return true
}

// 往服务中心注册自己前的操作
func (self *GoodlogProvider) BeforeInit(c core.Container) error {
	self.level = getLogLevel(c)
	return nil
}

func (sp *GoodlogProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (self *GoodlogProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		// 这里需要将参数展开，将配置注入到日志类，例如日志开关等
		c := params[0].(core.Container)
		if goodlogSvc != nil {
			return goodlogSvc, nil
		}
		goodlogSvc = &GoodlogService{c: c, level: self.level}
		return goodlogSvc, nil
	}

}

func (*GoodlogProvider) AfterInit(instance any) error {
	return nil
}

func getLogLevel(contain core.Container) byte {
	configSvs := contain.NewSingle(config.Name).(config.Service)
	config := configSvs.GetGoodLog()
	switch strings.ToLower(config.Level) {
	case "trace":
		return 0
	case "debug":
		return 1
	case "info":
		return 2
	case "warn":
		return 3
	case "error":
		return 4
	case "fatal":
		return 5
	case "off":
		return 6
	default:
		return 0
	}
}
