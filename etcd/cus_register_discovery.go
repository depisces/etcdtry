package etcd

import (
	"context"
	"fmt"
	"log"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type serviceCache struct {
	data map[string]string
	sync.RWMutex
}

var cache *serviceCache

func init() {
	cache = &serviceCache{
		data: make(map[string]string, 0),
	}
}

func getKey(serviceName string) string {
	return serviceName
}

// 服务注册 server端
func CusServiceRegister(serviceName, addr string) error {
	cli, err := GetEtcdClient()
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return err
	}
	key := getKey(serviceName)

	ctx := context.Background()
	//创建租约
	leaseRes, err := cli.Grant(ctx, 10)

	if err != nil {
		return err
	}
	//向etcd写数据
	_, err = cli.Put(ctx, key, addr, clientv3.WithLease(leaseRes.ID))
	if err != nil {
		return err
	}
	//服务端向etcd发送keepalive
	keepaliveCh, err := cli.KeepAlive(ctx, leaseRes.ID)
	if err != nil {
		return err
	}
	go func() {
		for item := range keepaliveCh {
			fmt.Printf("leaseID:%x,ttl:%v\n", item.ID, item.TTL)
		}
	}()
	return nil
}

// 服务发现 client端
func CusServiceDiscovery(serviceName string) string {
	cache.RLock()
	defer cache.RUnlock()
	return cache.data[serviceName]
}

// 第一次获取服务信息，监听key(地址)变化
func CusLoadService(serviceName string) {
	cli, _ := GetEtcdClient()
	ctx := context.Background()
	key := getKey(serviceName)
	getRes, err := cli.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
	if getRes.Count > 0 {
		cache.Lock()
		defer cache.Unlock()
		for _, item := range getRes.Kvs {
			cache.data[string(item.Key)] = string(item.Value)
		}
	}
}

func CusWatchService(serviceName string) {
	cli, _ := GetEtcdClient()
	ctx := context.Background()
	key := getKey(serviceName)
	rch := cli.Watch(ctx, key)
	for wres := range rch {
		for _, event := range wres.Events {
			if event.Type == clientv3.EventTypeDelete { //如果是删除事件
				cache.Lock()
				defer cache.Unlock()
				delete(cache.data, key)
				continue
			}

			if event.Type == clientv3.EventTypePut {
				cache.Lock()
				defer cache.Unlock()
				cache.data[key]=string(event.Kv.Value)
				continue

			}

		}
	}
}
