package main

import (
	"context"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"

	pb "github.com/LiangNing7/go-example/proto/09-middleware/proto"
	"github.com/LiangNing7/go-example/proto/09-middleware/server/middleware/auth"
	"github.com/LiangNing7/go-example/proto/09-middleware/server/middleware/cred"
	"github.com/LiangNing7/go-example/proto/09-middleware/server/middleware/recovery"
	"github.com/LiangNing7/go-example/proto/09-middleware/server/middleware/zap"
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

	// 新建 gRPC 服务器实例.
	grpcServer := grpc.NewServer(
		cred.TLSInterceptor(),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
			grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
			grpc_recovery.StreamServerInterceptor(recovery.RecoveryInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
			grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_recovery.UnaryServerInterceptor(recovery.RecoveryInterceptor()),
		)),
	)

	// 在 gRPC 服务器注册我们的服务.
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})
	log.Println(Address + " net.Listing with TLS and token...")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

// Route 实现 Route 方法.
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := &pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return res, nil
}
