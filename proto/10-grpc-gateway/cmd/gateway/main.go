package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/LiangNing7/go-example/proto/10-grpc-gateway/proto"
	"github.com/LiangNing7/go-example/proto/10-grpc-gateway/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// gRPC 服务地址.
	grpcAddr := ":9090"
	// HTTP Gateway 地址.
	httpAddr := ":8080"

	// 启动 gRPC 服务（在单独的 goroutine 中启动）.
	go func() {
		if err := server.RunGRPCServer(grpcAddr); err != nil {
			log.Fatalf("grpcServer.Serve err: %v", err)
		}
	}()

	// 创建 gRPC-Gateway 的 multiplexer.
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 注册 HTTP 转发.
	err := pb.RegisterSimpleHandlerFromEndpoint(context.Background(), mux, grpcAddr, opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
	log.Printf("HTTP gateway listening at %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("failed to serve HTTP gateway: %v", err)
	}
}
