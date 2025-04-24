package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/LiangNing7/go-example/proto/08-token/client/auth"
	pb "github.com/LiangNing7/go-example/proto/08-token/proto"
)

// Address 连接地址.
const Address string = ":8000"

var grpcClient pb.SimpleClient

func main() {
	// 从输入的证书文件中为客户端构造 TLS 凭证.
	creds, err := credentials.NewClientTLSFromFile("../key/test.pem", "blog.liangning7.cn")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	// 构建 Token.
	token := auth.Token{
		AppID:     "grpc_token",
		AppSecret: "123456",
	}

	// 连接服务器.
	conn, err := grpc.NewClient(
		Address,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&token),
	)
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立 gRPC 连接.
	grpcClient = pb.NewSimpleClient(conn)
	route()
}

// route 调用服务端Route方法
func route() {
	// 创建发送结构体
	req := pb.SimpleRequest{
		Data: "grpc",
	}
	// 调用我们的服务(Route方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := grpcClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res)
}
