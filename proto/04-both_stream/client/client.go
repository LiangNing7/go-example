package main

import (
	"context"
	"io"
	"log"
	"strconv"

	pb "github.com/LiangNing7/go-example/proto/04-both_stream/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Address 连接地址.
const Address string = ":8000"

var streamClient pb.StreamClient

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

	// 建立 gRPC 连接.
	streamClient = pb.NewStreamClient(conn)
	route()
	conversations()
}

// route 调用服务端Route方法
func route() {
	// 创建发送结构体
	req := pb.SimpleRequest{
		Data: "grpc",
	}
	// 调用我们的服务(Route方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := streamClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res.Value)
}

// conversations 调用服务端的 Conversations 方法.
func conversations() {
	// 调用服务端的 Conversations 方法，获取流.
	stream, err := streamClient.Conversations(context.Background())
	if err != nil {
		log.Fatalf("get conversations stream err: %v", err)
	}

	for n := range 5 {
		err := stream.Send(&pb.StreamRequest{
			Question: "stream client rpc " + strconv.Itoa(n),
		})
		if err != nil {
			log.Fatalf("stream request err: %v", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}

		// 打印返回值.
		log.Println(res.Answer)
	}
	// 最后关闭流.
	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("Conversations close stream err: %v", err)
	}
}
