package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/learn-grpc-03/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	creds, _ := credentials.NewClientTLSFromFile(
		"/root/workspace/learn-grpc-02/key/test.pem",
		"*.heliu.site",
	)

	var opts []grpc.DialOption
	// 不带TLS这里是grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// 连接server端，使用ssl加密通信
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	fmt.Printf("now-Time: %s\n", time.Now().Format(time.DateTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.UpdateOrders(ctx)

	if err != nil {
		log.Fatalf("error when calling SearchOrders: %v", err)
	}

	if err := stream.Send(&pb.Order{
		Id:          "1",
		Items:       []string{"6", "5", "4", "3", "2", "1"},
		Destination: "101",
	}); err != nil {
		log.Fatalf("send err: %v", err)
	}

	if err := stream.Send(&pb.Order{
		Id:          "2",
		Items:       []string{"61", "51", "41", "31", "21", "11"},
		Destination: "102",
	}); err != nil {
		log.Fatalf("send err: %v", err)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("recv err: %v", err)
	}

	log.Printf("Update Orders Res: %s\n", res.GetValue())
}
