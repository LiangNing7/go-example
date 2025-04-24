package main

import (
	"context"
	"io"
	"log"
	"strconv"

	pb "github.com/LiangNing7/go-example/proto/03-client_stream/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Address 连接地址.
const Address string = ":8000"

var streamClient pb.StreamClientClient

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
	streamClient = pb.NewStreamClientClient(conn)
	route()
	routeList()
}

// route 调用服务端 Route 方法.
func route() {
	// 创建发送结构体.
	req := pb.SimpleRequest{
		Data: "grpc",
	}
	// 调用我们的服务(Route方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变 RPC 的行为，比如超时/取消一个正在运行的 RPC。
	res, err := streamClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值.
	log.Println(res)
}

// routeList 调用服务端 RouteList 方法.
func routeList() {
	// 调用服务端的 RouteList 方法.
	stream, err := streamClient.RouteList(context.Background())
	if err != nil {
		log.Fatalf("Upload list err: %v", err)
	}

	for n := range 5 {
		// 向流中发送消息.
		err := stream.Send(&pb.StreamRequest{
			StreamData: "stream client rpc " + strconv.Itoa(n),
		})
		log.Println("StreamData: " + strconv.Itoa(n))
		// 发送也要检测 EOF，当服务端在消息没接收完前主动调用 SendAndClose() 关闭 stream，
		// 此时客户端还执行 Send()，则会返回 EOF 错误，所以这里需要加上 io.EOF 判断.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream request err: %v", err)
		}
	}

	// 关闭流并获取返回的消息.
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("RouteList get response err: %v", err)
	}
	log.Println(res)
}
