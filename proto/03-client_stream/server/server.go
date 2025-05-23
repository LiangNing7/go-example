package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "github.com/LiangNing7/go-example/proto/03-client_stream/proto"
	"google.golang.org/grpc"
)

// SimpleService 定义我们的服务.
type SimpleService struct {
	pb.UnimplementedStreamClientServer
}

// Route 实现 Route 方法.
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}

// RouteList 实现 RouteList 方法.
func (s *SimpleService) RouteList(srv pb.StreamClient_RouteListServer) error {
	for {
		// 从流中获取消息.
		res, err := srv.Recv()
		if err == io.EOF {
			// 发送结果，并关闭.
			return srv.SendAndClose(&pb.SimpleResponse{Value: "ok"})
		}
		if err != nil {
			return err
		}
		log.Println(res.StreamData)
	}
}

const (
	// Address 监听地址
	Address string = ":8000"
	// Network 网络通信协议
	Network string = "tcp"
)

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterStreamClientServer(grpcServer, &SimpleService{})

	// 用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
