package container

import (
	"errors"
	"fmt"
	"sync"
)

// 服务中心，也就是服务的容器，提供绑定服务和获取服务的能力
// 注册是全局一次性的，获取服务是业务场景中频繁获发生的，所以将服务中心的引用（服务容器）注入到 context
// 跟 http 路由器一样，服务中心是为服务提供注册、初始化、路由功能
type Container interface {
	// 绑定一个服务提供者，如果服务名称已经存在，会进行替换操作，返回 error
	// 服务注册其实跟 http 路由一样的原理
	Bind(provider ServiceProvider) error

	// 判断服务名称是否已经绑定服务提供者，防止重复绑定
	IsBind(name string) bool

	// 根据服务名称获取一个单例服务，如果这个服务名称未绑定服务提供者，那么会 panic，
	NewSingle(name string) interface{}

	// 根据服务名称获取一个非单例服务
	// 它是根据服务提供者注册的启动函数和传递的 params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	NewInstance(name string, params []interface{}) interface{}
}

// 服务容器的具体实现，其功能就是保存所有注册进来的服务和获取一个服务
type ServicesContainer struct {
	Container                            // 显示要求 ServicesContainer 实现 Container 接口
	providers map[string]ServiceProvider // 存储注册进来的服务提供者数据结构，这样就不用 import 服务提供者了，key 为服务名称
	instances map[string]interface{}     // 根据服务提供者数据结构实例化出具体实例并保存起来，key 为服务名称
	// 当 Bind 的时候有读有写，需要使用一个机制来保证 ServicesContainer 的并发性
	// ServicesContainer 是读多于写的数据结构，即 Bind 是一次性的，Make 是频繁的，用读写锁的性能优于互斥锁。
	lock sync.RWMutex
}

// 创建一个服务容器
func New() *ServicesContainer {
	container := &ServicesContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
	return container
}

// 将服务提供者的数据结构绑定到服务中心
func (self *ServicesContainer) Bind(provider ServiceProvider) error {
	name := provider.Name()
	if self.IsBind(name) {
		panic("Duplicate services provider name")
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	self.providers[name] = provider
	if provider.InitOnBoot() == true {
		if err := provider.BeforeInit(self); err != nil {
			return err
		}
		// 实例化服务提供者
		params := provider.Params(self)
		method := provider.RegisterProviderInstance(self)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		self.instances[name] = instance
		err = provider.AfterInit(instance)
		if err != nil {
			return err
		}
	}
	return nil
}

func (engine *ServicesContainer) IsBind(key string) bool {
	return engine.findServiceProvider(key) != nil
}

// 实例化一个单例服务
func (self *ServicesContainer) NewSingle(name string) interface{} {
	serv, err := self.make(name, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

// 不走单例，来一个实例化一个
func (self *ServicesContainer) NewInstance(key string, params []interface{}) interface{} {
	svc, err := self.make(key, params, true)
	if err != nil {
		panic(err)
	}
	return svc
}

// 实例化一个服务
func (self *ServicesContainer) make(name string, params []interface{}, forceNew bool) (interface{}, error) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	sp := self.findServiceProvider(name)
	if sp == nil {
		return nil, errors.New("provider structure " + name + " have not register")
	}
	// 强制实例化，也就是不走单例
	if forceNew {
		return self.makeNewInstance(sp, params)
	}
	// 单例获取，如果容器中已经有实例了就直接使用之
	if instance, ok := self.instances[name]; ok {
		return instance, nil
	}
	inst, err := self.makeNewInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	self.instances[name] = inst
	return inst, nil
}

// 按断是否有服务提供者数据结构，如果服务提供者没有把自己的数据结构传递到容器那还实例化个锤子
func (self *ServicesContainer) findServiceProvider(name string) ServiceProvider {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if sp, ok := self.providers[name]; ok {
		return sp
	}
	return nil
}

func (this *ServicesContainer) makeNewInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.BeforeInit(this); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(this)
	}
	method := sp.RegisterProviderInstance(this)
	ins, err := method(params...)
	err = sp.AfterInit(ins)
	if err != nil {
		return ins, err
	}
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

// 输出服务容器中注册的关键字
func (self *ServicesContainer) PrintProviders() []string {
	ret := []string{}
	for _, provider := range self.providers {
		name := provider.Name()
		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	//fmt.Println(ret)
	return ret
}

//func (hade *HadeContainer) Make(key string) (interface{}, error) {
//	return hade.make(key, nil, false)
//}
