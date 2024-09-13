package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"

	pb "example.com/learn-grpc-03/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	// 方法一
	// creds, err1 := credentials.NewServerTLSFromFile(
	//		"/root/workspace/learn-grpc/key/test.pem",
	//		"/root/workspace/learn-grpc/key/test.key",
	//	)
	//
	//	if err1 != nil {
	//		fmt.Printf("证书错误：%v", err1)
	//		return
	//	}

	// 方法二
	cert, err := tls.LoadX509KeyPair(
		"/root/workspace/learn-grpc-02/key/test.pem",
		"/root/workspace/learn-grpc-02/key/test.key")
	if err != nil {
		fmt.Printf("私钥错误：%v", err)
		return
	}
	creds := credentials.NewServerTLSFromCert(&cert)

	listen, _ := net.Listen("tcp", ":9090")
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterOrderManagementServer(grpcServer, &service{})

	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var orders = make(map[string]pb.Order, 8)

func init() {
	// 测试数据
	orders["1"] = pb.Order{Id: "1", Items: []string{"1", "2", "3", "4", "5", "7"}, Destination: "101"}
	orders["2"] = pb.Order{Id: "2", Items: []string{"6", "5", "4", "3", "2", "0"}, Destination: "102"}
}

var _ pb.OrderManagementServer = (*service)(nil)

type service struct {
	pb.UnimplementedOrderManagementServer
}

// UpdateOrders updates
// 在这段程序中，我们对每一个 Recv 都进行了处理
// 当发现 io.EOF (流关闭) 后，需要将最终的响应结果发送给客户端，同时关闭正在另外一侧等待的 Recv
func (s *service) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {
	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// Finished reading the order stream.
			return stream.SendAndClose(
				&wrapperspb.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order
		orders[order.Id] = *order

		log.Println("Order ID ", order.Id, ": Updated")
		ordersStr += order.Id + ", "
	}
}
