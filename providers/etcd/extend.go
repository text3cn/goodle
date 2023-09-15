package etcd

import (
	"context"
	"fmt"
	"github.com/text3cn/goodle/config"
	"github.com/text3cn/goodle/providers/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func client() (*clientv3.Client, config.EtcdConfig) {
	cfg := config.Instance().GetDiscovery()
	// 创建 etcd 客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Server.Endpoints,
		DialTimeout: time.Duration(cfg.Server.DialTimeoutSecods) * time.Second,
	})
	if err != nil {
		logger.Pink("client() ERR:", err)
	}
	return cli, cfg
}

// 启动心跳机制
func sendHeartbeats(ctx context.Context, cli *clientv3.Client, id clientv3.LeaseID, interval byte) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	defer cli.Close()

	for {
		select {
		case <-ticker.C:
			// 定时发送心跳
			_, err := cli.KeepAliveOnce(context.TODO(), id)
			if err != nil {
				log.Printf("Send heartbeats error：%v\n", err)
			}
		case <-ctx.Done():
			// 程序退出时停止心跳
			fmt.Println("心跳停止")
			return
		}
	}
}
