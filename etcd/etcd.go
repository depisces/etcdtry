package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetEtcdEndpoints() []string {
	return []string{"http://127.0.0.1:2379"}
}

func GetEtcdClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: GetEtcdEndpoints(),
	})
	return cli, err
}
