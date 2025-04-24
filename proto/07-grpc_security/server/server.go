package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/LiangNing7/go-example/proto/07-grpc_security/proto"
)

// SimpleService 定义我们的服务.
type SimpleService struct {
	pb.UnimplementedSimpleServer
}

const (
	// Address 监听地址.
	Address string = ":8000"
	// Network 网络通信协议.
	Network string = "tcp"
)

func main() {
	// 监听本地端口.
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	// 从输入证书文件和密钥文件为服务端构建 TSL 凭证.
	creds, err := credentials.NewServerTLSFromFile("../key/test.pem", "../key/test.key")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	// 新建 gRPC 实例，并开启 TLS 认证.
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// 在 gRPC 服务器中注册我们的服务.
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	log.Println(Address + " net.Listing with TLS and token...")
	// 用服务器 Server() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用.
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}
