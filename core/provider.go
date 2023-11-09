// 一切皆服务，制定向服务中心注册服务时需要实现的方法标准
package core

// 定义如何创建一个新实例，所有服务容器的创建服务
type NewInstanceFunc func(...interface{}) (interface{}, error)

// 定义服务提供者需要实现的接口
type ServiceProvider interface {
	// 服务名称
	Name() string

	// 决定是否在注册（程序启动）时实例化这个服务，true 则  bind() 时直接实例化，false 为手动获取实例的时候才进行实例化
	InitOnBoot() bool

	// 服务调用者，实例化一个服务时传递的参数
	Params(Container) []interface{}

	// 实例化一个服务提供者，并保存起来，这样服务中心不需要 import 就持有了各服务的实例
	// 函数在 Golang 中是一等公民，各个服务提供者通过将他们实例创建的方法通过回调函数的形式传递过来，
	// 这样服务中心就在不需要 import 各个服务文件的情况下持有了服务的实力，
	// 然后服务中心会被注入到 context 中，那么就可以在任何有 context 的地方调用任何服务了
	RegisterProviderInstance(Container) NewInstanceFunc

	// 实例化服务前的初始化操作，如果有 error 则不做服务实例化
	BeforeInit(Container) error

	// 实例化服务后的操作，将实例传递过去做一些服务初始化操作
	AfterInit(any) error
}
