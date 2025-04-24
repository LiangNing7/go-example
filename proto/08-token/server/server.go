package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/LiangNing7/go-example/proto/08-token/proto"
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

	// 普通方法：一元拦截器(grpc.UnaryServerInterceptor)
	var interceptor grpc.UnaryServerInterceptor = func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 拦截普通方法请求，验证 Token.
		err = Check(ctx)
		if err != nil {
			return
		}
		// 继续处理请求.
		return handler(ctx, req)
	}

	// 新建 gRPC 实例，并开启 TLS 认证.
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(interceptor))
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
	// 添加拦截器后，方法里省略 Token 认证.
	// // 检测 Token 是否有效.
	// if err := Check(ctx); err != nil {
	// 	return nil,err
	// }
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}

// Check 验证 Token.
func Check(ctx context.Context) error {
	// 从上下文获取元数据.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "获取 Token 失败")
	}

	var (
		appID     string
		appSecret string
	)

	if value, ok := md["app_id"]; ok {
		appID = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appID != "grpc_token" || appSecret != "123456" {
		return status.Errorf(codes.Unauthenticated, "Token 无效: app_id=%s, app_secret=%s", appID, appSecret)
	}
	return nil
}
