// 一切皆服务，制定向服务中心注册服务时需要实现的方法标准
package container

// 定义如何创建一个新实例，所有服务容器的创建服务
type NewInstance func(...interface{}) (interface{}, error)

// 定义服务提供者需要实现的接口
type ServiceProvider interface {
	// 服务名称
	Name() string

	// 实例化一个服务提供者，并保存起来，
	// 函数在 Golang 中是一等公民，各个服务提供者通过将他们实例创建的方法通过回调函数的形式传递过来，
	// 这样服务中心就在不需要 import 各个服务文件的情况下持有了服务的实力，
	// 然后服务中心会被注入到 context 中，那么就可以在任何有 context 的地方调用任何服务了
	RegisterProviderInstance(Container) NewInstance

	// 决定是否在注册（程序启动）时实例化这个服务，false 代表在第一次 new 的时候进行实例化，
	// 也就是用到的时候再实例化，没有触发使用就不实例化，也就是按需加载
	IsDefer() bool

	// 定义传递给 NewInstance 的参数，可以自定义多个参数
	Params(Container) []interface{}

	// 实例化服务前的初始化操作，如果有 error 则整个服务实例化就会实例化失败
	Boot(Container) error
}
