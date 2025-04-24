package main

import (
	"context"
	"log"
	"time"

	pb "github.com/LiangNing7/go-example/proto/06-deadline/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Address 连接地址.
const Address string = ":8000"

var grpcClient pb.SimpleClient

func main() {
	// 连接服务器.
	conn, err := grpc.NewClient(
		Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	// 建立 gRPC 连接.
	grpcClient = pb.NewSimpleClient(conn)
	route(ctx, 2)
}

// route 调用服务端 Route 方法.
func route(ctx context.Context, deadlines time.Duration) {
	// 设置 3s 超时时间.
	clientDeadline := time.Now().Add(time.Duration(deadlines * time.Second))
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()
	// 创建发送结构体.
	req := pb.SimpleRequest{
		Data: "grpc",
	}
	// 调用我们的服务(Route 方法).
	// 传入超时时间为 3s 的 ctx.
	res, err := grpcClient.Route(ctx, &req)
	if err != nil {
		// 获取错误状态.
		statu, ok := status.FromError(err)
		if ok {
			// 判断是否为调用超时.
			if statu.Code() == codes.DeadlineExceeded {
				log.Fatalln("Route timeout!")
			}
		}
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值.
	log.Println(res.Value)
}
