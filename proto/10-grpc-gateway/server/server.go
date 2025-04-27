package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/LiangNing7/go-example/proto/10-grpc-gateway/proto"
)

// SimpleService 定义我们的服务
type SimpleService struct {
	pb.UnimplementedSimpleServer
}

func RunGRPCServer(addr string) error {
	// 监听本地端口
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(addr + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	// 用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	return err
}

// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}
