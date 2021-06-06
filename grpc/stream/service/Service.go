package main

import (
	"github.com/hxzhouh/go_example/grpc/stream/service/internel"
	"github.com/hxzhouh/go_example/grpc/stream/stream"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	log.Println("app running", zap.String("time now", time.Now().Format("2006-01-02 15:04:05")))
	lis, err := net.Listen("tcp", "localhost:50001") //开启监听
	if err != nil {
		log.Fatal("failed to listen: %v")
	}
	s := grpc.NewServer()                                             //新建一个grpc服务
	stream.RegisterEchoStreamServer(s, &internel.StreamServiceImpl{}) //这个服务和上述的服务结构联系起来，这样你新建的这个服务里面就有那些类型的方法
	if err := s.Serve(lis); err != nil {                              //这个服务和你的监听联系起来，这样外界才能访问到啊
		log.Fatal("failed to serve:", zap.Error(err))
	}
}
