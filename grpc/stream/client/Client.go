package main

import (
	"context"
	"github.com/hxzhouh/go_example/grpc/stream/stream"
	"google.golang.org/grpc"
	"log"
	"time"
)

const PORT = "localhost:50001"

func main() {
	conn, err := grpc.Dial(PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := stream.NewEchoStreamClient(conn)
	err = SayRoute(client)
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}
}

func SayRoute(client stream.EchoStreamClient) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	s, _ := client.SayHelloStream(ctx)
	for i := 0; i < 10; i++ {
		s.Send(&stream.StreamRequest{
			MessageType: 1,
			Body: &stream.StreamRequestBody{
				MessageId:  0,
				ReceiverId: 0,
				Name:       "hello,world",
			},
		})
		time.Sleep(100 * time.Microsecond)
	}
	time.Sleep(10 * time.Second)
	_ = s.CloseSend()
	return nil
}
