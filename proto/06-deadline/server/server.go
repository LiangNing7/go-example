package main

import (
	"context"
	"log"
	"net"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/LiangNing7/go-example/proto/06-deadline/proto"
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
	log.Println(Address + " net.Listing...")
	// 新建 gRPC 实例.
	grpcServer := grpc.NewServer()
	// 在 gRPC 服务器中注册我们的服务.
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	// 用服务器 Server() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用.
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

// Route 实现 Route 方法.
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	data := make(chan *pb.SimpleResponse, 1)
	go handle(ctx, req, data)
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Client canceled, abandoning.")
	}
}

func handle(ctx context.Context, req *pb.SimpleRequest, data chan<- *pb.SimpleResponse) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		runtime.Goexit() // 超时后退出该 goroutine.
	case <-time.After(4 * time.Second):
		res := &pb.SimpleResponse{
			Code:  200,
			Value: "hello " + req.Data,
		}
		// // 修改数据库前进行超时判断.
		// if ctx.Err() == context.Canceled {
		// 	...
		// 	// 如果已经超时，则退出.
		// }
		data <- res
	}
}
