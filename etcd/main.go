package main

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"log"
	"time"
)

var client *clientv3.Client

func main() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.209.128:12379", "192.168.209.128:22379", "192.168.209.128:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	if client == nil {
		log.Fatal("client is nil")
	}
	timeout := time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_, err = client.Put(ctx, "sample_key", "sample_value")
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}

	resp, err := client.Get(ctx, "sample_key")
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	log.Println(resp.Kvs[0].String())
	cancel()
	defer client.Close()
}
