package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	pb "github.com/LiangNing7/go-example/proto/04-both_stream/proto"
)

// StreamService 定义我们的服务
type StreamService struct {
	pb.UnimplementedStreamServer
}

const (
	// Address 监听地址.
	Address string = ":8000"
	// Network 通络通信协议.
	Network string = "tcp"
)

func main() {
	// 监听本地端口.
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	log.Println(Address + " net.Listing...")
	// 新建 gRPC 服务器实例.
	grpcServer := grpc.NewServer()
	// 在 gRPC 服务器中注册我们的服务.
	pb.RegisterStreamServer(grpcServer, &StreamService{})

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

// Conversations 实现 Conversations 方法.
func (s *StreamService) Conversations(srv pb.Stream_ConversationsServer) error {
	n := 1
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = srv.Send(&pb.StreamResponse{
			Answer: "from stream server answer: the " + strconv.Itoa(n) + " question is " + req.Question,
		})
		if err != nil {
			return err
		}
		n++
		log.Printf("from stream client question: %s", req.Question)
	}
}
