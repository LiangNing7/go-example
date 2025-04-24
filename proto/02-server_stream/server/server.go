package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"

	pb "github.com/LiangNing7/go-example/proto/02-server_stream/proto"
)

// StreamService 定义我们的服务.
type StreamService struct {
	pb.UnimplementedStreamServerServer
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
	log.Println(Address, " net.Listing...")
	// 新建 gRPC 服务器实例.
	// 默认单次接收最大消息长度为 `1024 * 1024 * 4`bytes(4M)，单次发送消息最大长度为 `math.MaxInt32`bytes.
	// grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*4), grpc.MaxSendMsgSize(math.MaxInt32))
	grpcServer := grpc.NewServer()
	// 在 gRPC 服务器中注册我们的服务.
	pb.RegisterStreamServerServer(grpcServer, &StreamService{})

	// 用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用.
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

// Route 实现 Route 方法.
func (s *StreamService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}

// // ListValue 实现 ListValue 方法.
// func (s *StreamService) ListValue(req *pb.SimpleRequest, srv pb.StreamServer_ListValueServer) error {
// 	for n := range 5 {
// 		// 向流中发送消息，默认每次 send 消息的最大长度为 `math.MaxInt32`bytes
// 		err := srv.Send(&pb.StreamResponse{
// 			StreamValue: req.Data + strconv.Itoa(n),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// ListValue 实现 ListValue 方法.
func (s *StreamService) ListValue(req *pb.SimpleRequest, srv pb.StreamServer_ListValueServer) error {
	for n := range 15 {
		// 向流中发送消息，默认每次 send 消息最大长度为`math.MaxInt32`bytes.
		err := srv.Send(&pb.StreamResponse{
			StreamValue: req.Data + strconv.Itoa(n),
		})
		if err != nil {
			return err
		}
		log.Println(n)
		time.Sleep(1 * time.Second)
	}
	return nil
}
