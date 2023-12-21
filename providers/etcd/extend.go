package etcd

import (
	"context"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/goodlog"
	"github.com/text3cn/goodle/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func client() (*clientv3.Client, types.EtcdConfig) {
	cfg := config.Instance().GetEtcd()
	// 创建 etcd 客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Server.Endpoints,
		DialTimeout: time.Duration(cfg.Server.DialTimeoutSecods) * time.Second,
	})
	if err != nil {
		goodlog.Info("client() ERR:", err)
	}
	return cli, cfg
}

// 启动心跳机制
func sendHeartbeats(self *etcdService, cli *clientv3.Client, id clientv3.LeaseID, interval byte) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer self.ServiceRegister() // 重新注册
	defer ticker.Stop()
	defer cli.Close()

	for {
		select {
		case <-ticker.C:
			// 定时发送心跳
			_, err := cli.KeepAliveOnce(context.TODO(), id)
			if err != nil {
				log.Printf("[Etcd 续约失败] Send heartbeats error：%v\n", err)
				return
			}
			time.Sleep(time.Second * 6)
		}
	}
}
