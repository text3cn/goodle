package etcd

import (
	"context"
	"fmt"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/goodlog"
	"github.com/text3cn/goodle/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Service interface {
}

type etcdService struct {
	Service
	c core.Container
}

func Instance() *etcdService {
	return etcdInstance
}

// 服务注册
func (self *etcdService) ServiceRegister() {
	cli, cfg := client()
	// 设置上下文和服务信息
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.DialTimeoutSecods)*time.Second)
	defer cancel()

	// 分配租约时间为 30 秒
	leaseResp, err := cli.Grant(ctx, 30)
	if err != nil {
		goodlog.Pink(leaseResp.Error)
		return
	}

	serviceName := cfg.Client.ServiceName
	serviceAddr := cfg.Client.ServiceAddr

	// 注册服务到 etcd
	key := fmt.Sprintf("/services/%s", serviceName)
	_, err = cli.Put(ctx, key, serviceAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		goodlog.Error(err)
		return
	}

	// 启动心跳每 20 秒续约一次
	go sendHeartbeats(self, cli, leaseResp.ID, 20)

	goodlog.Trace("注册服务 " + serviceName + " 到 etcd 成功")
}

// 服务下线
func (self *etcdService) ServiceOffline() {
	cli, cfg := client()
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.DialTimeoutSecods)*time.Second)
	defer cancel()

	// 删除服务的注册信息
	key := fmt.Sprintf("/services/%s", cfg.Client.ServiceName)

	_, err := cli.Delete(ctx, key)
	if err != nil {
		goodlog.Error(err)
	}

	goodlog.Pink("服务 " + cfg.Client.ServiceName + " 已下线并从 etcd 中注销")
}

// 发现服务
func (self *etcdService) ServiceDiscovery(serviceName string, callback func(services []string)) {
	cli, cfg := client()
	defer cli.Close()
	var token = "none" // 将所有节点拼接为 token 用于判断是否有节点变更
	discoverService := func(cli *clientv3.Client, serviceName string) {
		services := self.GetServices(serviceName)
		_token := ""
		for _, v := range services.List {
			_token += v
		}

		// 没有节点变更就不用让 Etcd 客户端感知了
		if token == "none" || _token != token {
			goodlog.Trace("发现服务 "+serviceName+" 地址列表为 ", services)
			callback(services.List)
			if len(services.List) == 0 {
				goodlog.Pink("none services " + serviceName)
			}
		}
		token = _token
	}

	// 通过服务名进行服务发现
	discoverService(cli, serviceName)

	// 启动协程不停地去发现服务
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.DiscoveryIntervalSeconds) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				discoverService(cli, serviceName)
			}
		}
		defer func() {
			if err := recover(); err != nil {
				goodlog.Error(err)
			}
		}()
	}()
}

// 根据服务名称列出所有服务
func (self *etcdService) GetServices(serviceName string) (ret types.GetServicesDTO) {
	cli, cfg := client()
	defer cli.Close()
	dialTimeout := time.Duration(cfg.Server.DialTimeoutSecods) * time.Second

	// 设置上下文和服务信息
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	// 查询以服务信息前缀开头的键值对
	servicePrefix := "/services/" + serviceName
	resp, err := cli.Get(ctx, servicePrefix, clientv3.WithPrefix())
	if err != nil {
		goodlog.Pink(err.Error())
	}
	list := []string{}
	for _, kv := range resp.Kvs {
		serviceAddress := string(kv.Value)
		list = append(list, serviceAddress)
	}
	ret.List = list
	return
}

// 根据服务名称列出所有服务
func (self *etcdService) ServicesList() (ret map[string][]string) {
	cli, cfg := client()
	defer cli.Close()
	dialTimeout := time.Duration(cfg.Server.DialTimeoutSecods) * time.Second
	// 设置上下文和服务信息
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	// 查询以服务信息前缀开头的键值对
	servicePrefix := "/services/"
	resp, err := cli.Get(ctx, servicePrefix, clientv3.WithPrefix())
	if err != nil {
		goodlog.Pink(err.Error())
	}
	ret = make(map[string][]string)
	for _, kv := range resp.Kvs {
		serviceName := string(kv.Key[len(servicePrefix):])
		serviceAddress := string(kv.Value)
		ret[serviceName] = append(ret[serviceName], serviceAddress)
	}
	return
}
