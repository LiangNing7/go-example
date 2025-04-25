# gRPC 介绍及环境配置

## RPC 介绍

gRPC 是 RPC 协议的 Go 语言实现，因此在介绍 gRPC 之前，有必要先了解 RPC 协议的基本概念。

根据维基百科的定义，RPC（Remote Procedure Call，远程过程调用）是一种计算机通信协议。该协议允许运行在一台计算机上的程序调用另一台计算机上的子程序，而开发者无需为这种交互编写额外的代码。

通俗来说，RPC 的核心思想是：服务端实现一个函数，客户端通过 RPC 框架提供的接口，可以像调用本地函数一样调用该函数，并获取返回值。RPC 屏蔽了底层的网络通信细节，使开发人员无需关注网络编程的复杂性，从而能够将更多精力投入到业务逻辑的实现中，大幅提升开发效率。

RPC 的调用过程如下图所示：

![image-20250423121607100](http://images.liangning7.cn/typora/202504231218288.png)

RPC 调用的具体流程如下：

1. 客户端通过本地调用的方式调用客户端存根（Client Stub）；
2. 客户端存根将参数打包（也称为 Marshalling）成一个消息，并发送该消息；
3. 客户端所在的操作系统（OS）将消息发送到服务端；
4. 服务端接收到消息后，将消息传递给服务端存根（Server Stub）；
5. 服务端存根将消息解包（也称为 Unmarshalling），得到参数；
6. 服务端存根调用服务端的子程序（函数），完成处理后，将结果按照相反的步骤返回给客户端。

需要注意的是，Stub 负责处理参数和返回值的序列化（Serialization）、参数的打包与解包，以及网络层的通信。在 RPC 中，客户端的 Stub 通常被称为“Stub”，而服务端的 Stub 通常被称为“Skeleton”。

> tips:
>
> Stub 是存根（代理的意思）

目前，业界有很多优秀的 RPC 协议，例如腾讯的 Tars、阿里的 Dubbo、微博的 Motan、Facebook 的 Thrift、RPCX 等。但使用最多的还是 gRPC。

## gRPC 介绍

![image-20250425121551933](http://images.liangning7.cn/typora/202504251215106.png)

gRPC 是由谷歌开发的一种高性能、开源且支持多种编程语言的通用 RPC 框架，基于 HTTP/2 协议开发，并默认采用 Protocol Buffers 作为数据序列化协议。gRPC 具有以下特性：

* **语言中立：**支持多种编程语言，例如 Go、Java、C、C++、C#、Node.js、PHP、Python、Ruby 等；
* **基于 IDL 定义服务：**通过 IDL（Interface Definition Language）文件定义服务，并使用 proto3 工具生成指定语言的数据结构、服务端接口以及客户端存根。这种方法能够解耦服务端和客户端，实现客户端与服务端的并行开发；
* **基于 HTTP/2 协议：**通信协议基于标准的 HTTP/2 设计，支持双向流、消息头压缩、单 TCP 的多路复用以及服务端推送等能力；
* **支持 Protocol Buffer 序列化：**Protocol Buffer（简称 Protobuf）是一种与语言无关的高性能序列化框架，可以减少网络传输流量，提高通信效率。此外，Protobuf 语法简单且表达能力强，非常适合用于接口定义。

> tips:
>
> gRPC 的全称并非“golang Remote Procedure Call”，而是“google Remote Procedure Call”。

与许多其他 RPC 框架类似，gRPC 也通过 IDL 语言来定义接口（包括接口名称、传入参数和返回参数等）。在服务端，gRPC 服务实现了预定义的接口。在客户端，gRPC 存根提供了与服务端相同的方法。

### HTTP2

gRPC 在底层使用 HTTP/2。与 HTTP/1.1 相比，它具有：

* **多路复用**：在同一连接上并行处理多个请求；
* **支持双向流**：客户端和服务器可同时发送和接收数据；
* **内置头部压缩：**加快数据传输速度。

![image-20250425122154027](http://images.liangning7.cn/typora/202504251221180.png)

得益于此，gRPC 能在低延迟下搞定实时通信。

### gRPC vs REST

| 特性     | REST                | gRPC             |
| :------- | :------------------ | :--------------- |
| 数据格式 | JSON / XML          | Protocol Buffers |
| 流式     | 罕见 / 需自己造轮子 | 内置             |
| 传输层   | HTTP/1.1            | HTTP/2           |
| 性能     | 冗长                | 紧凑 & 快        |
| 开发体验 | 手动写文档          | 自动生成代码     |

## 安装 protobuf 

要安装 Protobuf 文件的编译器 protoc。protoc 需要 protoc-gen-go 插件和一些其他插件来完成 Go 语言的代码转换，因此我们需要安装 protoc 和一些其他工具。它们的安装方法比较简单，具体分为以下两步：

1. 安装 protoc 命令。

   安装命令如下：

   ```bash
   $ cd /tmp/
   $ wget https://github.com/protocolbuffers/protobuf/releases/download/v29.1/protoc-29.1-linux-x86_64.zip
   $ mkdir protobuf-29.1/
   $ unzip protoc-29.1-linux-x86_64.zip -d protobuf-29.1/
   $ cd protobuf-29.1/
   $ sudo cp -a include/* /usr/local/include/
   $ sudo cp bin/protoc /usr/local/bin/
   $ protoc --version # 查看 protoc 版本，成功输出版本号，说明安装成功
   libprotoc 29.1
   ```

   > 提示：这里我们安装的 protoc 版本是 29.1。如果你安装了其他版本，后续执行 protoc 命令报错，可能需要你根据所安装的版本进行命令参数适配。

2. 安装 Protocol Buffers 编译器插件。

   protoc 编译工具依赖一些编译插件来生成相应的 Go 代码。miniblog 项目所依赖的 Protocol Buffers 编译插件如下图所示。

   ![protoc-go-tools](http://images.liangning7.cn/typora/202504091157021.png)

   安装命令如下：

   ```bash
   $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
   $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
   $ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
   $ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0
   $ go install github.com/onexstack/protoc-gen-defaults@v0.0.2
   ```

   当你第一次执行 go install 命令的时候，因为本地无缓存，需要下载所有的依赖模块，所以安装速度会比较慢，请你耐心等待。

3. 安装相关包

   安装 golang 的proto工具包

   ```bash
   $ go get -u github.com/golang/protobuf/proto
   ```

   安装 gRPC 包

   ```bash
   $ go get -u google.golang.org/grpc
   ```

# gRPC 基础

gRPC主要有4种请求和响应模式，分别是简单模式(`Simple RPC`)、服务端流式(`Server-side streaming RPC`)、客户端流式(`Client-side streaming RPC`)、和双向流式(`Bidirectional streaming RPC`)。

* **简单模式（Simple RPC）：**这是最基本的 gRPC 调用形式。客户端发起一个请求，服务端返回一个响应。定义格式为

  ```protobuf
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  ```

* **服务端流模式（Server-side streaming RPC）：**客户端发送一个请求，服务端返回数据流，客户端从流中依次读取数据直到流结束。定义格式为

  ```protobuf
  rpc SayHello (HelloRequest) returns (stream HelloReply) {}
  ```

* **客户端流模式（Client-side streaming RPC）：**客户端以数据流的形式连续发送多条消息至服务端，服务端在处理完所有数据之后返回一次响应。定义格式为

  ```protobuf
  rpc SayHello (stream HelloRequest) returns (HelloReply) {}
  ```

* **双向数据流模式（Bidirectional streaming RPC）：**客户端和服务端可以同时以数据流的方式向对方发送消息，实现实时交互。定义格式为

  ```protobuf
  rpc SayHello (stream HelloRequest) returns (stream HelloReply) {}
  ```

## 简单模式 RPC

### 新建 proto 文件

> 如果你想先学习 .proto文件的编写，可以先跳转[这里](#proto 文件介绍)进行学习

主要是定义我们服务的方法以及数据格式：

1. 定义头部：

   ```protobuf
   syntax = "proto3"; // 协议为 proto3.                                                  
   package proto;
   
   option go_package = "./;proto";
   ```

2. 定义发送消息的信息：

   ```protobuf
   // 定义发送请求消息.
   message SimpleRequest{
     // 定义发送的参数，采用驼峰命令方式，小写加下划线.
     // 如：student_name.
     // 声明方式：参数类型 参数名 标识号（不可重复）
     // 标识符用于在编译后的二进制消息格式中对字段进行识别
     // 一旦 Protobuf 消息投入使用，字段的标识符就不应再修改。
     // 数字标签的取值范围为 `[1, 536870911]`，
     // 其中 19000 至 19999 为保留数字，不能使用。  
     string data = 1;                     
   }
   ```

3. 定义响应信息：

   ```protobuf
   // 定义响应消息.
   message SimpleResponse{
     // 定义接收的参数.
     // 参数类型 参数名 标识号（不可重复）
     int32 code = 1;
     string value = 2;
   }
   ```

4. 定义服务方法 Route：

   ```protobuf
   // 定义我们的服务（可定义多个服务,每个服务可定义多个接口）
   service Simple{
       rpc Route (SimpleRequest) returns (SimpleResponse){};
   }
   ```

完整代码位于：[simple.proto](https://github.com/LiangNing7/go-example/blob/main/proto/01-simple_proto/proto/simple.proto)

```protobuf
syntax = "proto3"; // 协议为 proto3.

package proto;

option go_package = "./;proto";

// 定义发送请求消息.
message SimpleRequest{
  // 定义发送的参数，采用驼峰命令方式，小写加下划线.
  // 如：student_name.
  // 声明方式：参数类型 参数名 标识号（不可重复）
  // 标识符用于在编译后的二进制消息格式中对字段进行识别。
  // 一旦 Protobuf 消息投入使用，字段的标识符就不应再修改。
  // 数字标签的取值范围为 `[1, 536870911]`，
  // 其中 19000 至 19999 为保留数字，不能使用。
  string data = 1;
}

// 定义响应消息.
message SimpleResponse{
  // 定义接收的参数.
  // 参数类型 参数名 标识号（不可重复）
  int32 code = 1;
  string value = 2;
}

// 定义我们的服务（可定义多个服务，每个服务课定义多个接口）.
service Simple{
  rpc Route (SimpleRequest) returns (SimpleResponse){}
}
```

最后编译 proto 文件：

进入 `simple.proto`文件所在目录，运行：

```bash
$ protoc --go_out=. --go-grpc_out=. ./simple.proto
```

编译之后，在该目录有如下三个文件：

```bash
$ ls
simple_grpc.pb.go  simple.pb.go  simple.proto
```

### 创建 Server 端

1. 定义我们的服务，并实现 Route 方法

   ```go
   import (
   	"context"
   	"log"
   	"net"
   
   	"google.golang.org/grpc"
   
   	pb "go-grpc-example/proto"
   )
   // SimpleService 定义我们的服务
   type SimpleService struct {
   	pb.UnimplementedSimpleServer
   }
   
   // Route 实现Route方法
   func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
   	res := pb.SimpleResponse{
   		Code:  200,
   		Value: "hello " + req.Data,
   	}
   	return &res, nil
   }
   ```

   该方法需要传入 RPC 的上下文 `context.Context`，它的作用结束`超时`或`取消`的请求。

2. 启动 gRPC 服务器

   ```go
   const (
   	// Address 监听地址
   	Address string = ":8000"
   	// Network 网络通信协议
   	Network string = "tcp"
   )
   
   func main() {
   	// 监听本地端口
   	listener, err := net.Listen(Network, Address)
   	if err != nil {
   		log.Fatalf("net.Listen err: %v", err)
   	}
   	log.Println(Address + " net.Listing...")
   	// 新建gRPC服务器实例
   	grpcServer := grpc.NewServer()
   	// 在gRPC服务器注册我们的服务
   	pb.RegisterSimpleServer(grpcServer, &SimpleService{})
   
   	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
   	err = grpcServer.Serve(listener)
   	if err != nil {
   		log.Fatalf("grpcServer.Serve err: %v", err)
   	}
   }
   ```

   里面每个方法的作用都有注释，这里就不解析了。

完整代码位于：[server.go](https://github.com/LiangNing7/go-example/blob/main/proto/01-simple_proto/server/server.go)

```go
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/LiangNing7/go-example/proto/01-simple_proto/proto"
)

// SimpleService 定义我们的服务
type SimpleService struct {
	pb.UnimplementedSimpleServer
}

const (
	// Address 监听地址
	Address string = ":8000"
	// Network 网络通信协议
	Network string = "tcp"
)

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	// 用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}
```

运行服务端：

```bash
$ cd go-example/proto/01-simple_proto/server
$ go run server.go
:8000 net.Listing...
```

### 创建 Client 端

完整代码位于：[client.go](https://github.com/LiangNing7/go-example/blob/main/proto/01-simple_proto/client/client.go)

```go
package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/LiangNing7/go-example/proto/01-simple_proto/proto"
)

// Address 连接地址
const Address string = ":8000"

var grpcClient pb.SimpleClient

func main() {
	// 连接服务器
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
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
```

运行客户端：

```bash
$ cd go-example/proto/01-simple_proto/client
$ go run client.go
code:200 value:"hello grpc"
```

成功调用 Server 端的 Route 方法并获取返回的数据。

## 服务端流式 RPC

当数据量大或者需要不断传输数据时候，我们应该使用流式 RPC，它允许我们边处理边传输数据。这里先介绍服务器流式 RPC。

`服务端流式RPC`：客户端发送请求到服务器，拿到一个流去读取返回的消息序列。 客户端读取返回的流，直到里面没有任何消息。

**情景模拟：实时获取股票走势。**

1. 客户端要获取某原油股的实时走势，客户端发送一个请求
2. 服务端实时返回该股票的走势

### 新建 proto 文件

主要是定义我们服务的方法以及数据格式：

1. 定义头部：

   ```protobuf
   syntax = "proto3";
   
   package proto;
   
   option go_package = ".;proto";
   ```

2. 定义发送消息的信息：

   ```protobuf
   // 定义发送请求消息.
   message SimpleRequest{
     string data = 1;
   }
   ```

3. 定义响应信息：

   ```protobuf
   // 定义响应信息.
   message SimpleResponse{
     int32 code = 1;
     string value = 2;  
   }
   
   // 定义流式响应信息.
   message StreamResponse {
     // 流式响应数据.
     string stream_value = 1;
   }
   ```

4. 定义服务方法 Route：

   ```protobuf
   // 定义我们的服务（可定义多个服务，每个服务可定义多个接口）.
   service StreamServer {
     rpc Route (SimpleRequest) returns (SimpleResponse) {};
   
     // 服务端流式 rpc，在响应数据前添加 stream.
     rpc ListValue(SimpleRequest) returns (stream StreamResponse){};
   }
   ```

完整代码位于：[server_stream.proto](https://github.com/LiangNing7/go-example/blob/main/proto/02-server_stream/proto/server_stream.proto)

最后编译 proto 文件：

进入 `server_stream.proto`文件所在目录，运行：

```bash
$ protoc --go_out=. --go-grpc_out=. ./server_stream.proto
```

### 创建 Server 端

1. 定义我们的服务，并实现 `StreamServerServer` 接口。

   ```go
   // StreamService 定义我们的服务.
   type StreamService struct {
   	pb.UnimplementedStreamServerServer
   }
   
   // Route 实现 Route 方法.
   func (s *StreamService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
   	res := pb.SimpleResponse{
   		Code:  200,
   		Value: "hello " + req.Data,
   	}
   	return &res, nil
   }
   
   // ListValue 实现 ListValue 方法.
   func (s *StreamService) ListValue(req *pb.SimpleRequest, srv pb.StreamServer_ListValueServer) error {
   	for n := range 5 {
   		// 向流中发送消息，默认每次 send 消息的最大长度为 `math.MaxInt32`bytes
   		err := srv.Send(&pb.StreamResponse{
   			StreamValue: req.Data + strconv.Itoa(n),
   		})
   		if err != nil {
   			return err
   		}
   	}
   	return nil
   }
   ```

   你可能觉得比较迷惑，`ListValue` 的参数和返回值是怎样确定的。其实这些都是编译 proto 时生成的文件中有定义，我们只需要实现就可以了。

   ```go
   // 在生成的 server_stream_grpc.pb.go 中有如下定义：
   
   // StreamServerServer is the server API for StreamServer service.
   // All implementations must embed UnimplementedStreamServerServer
   // for forward compatibility.
   //
   // 定义我们的服务（可定义多个服务，每个服务可定义多个接口）.
   type StreamServerServer interface {
   	Route(context.Context, *SimpleRequest) (*SimpleResponse, error)
   	// 服务端流式 rpc，在响应数据前添加 stream.
   	ListValue(*SimpleRequest, grpc.ServerStreamingServer[StreamResponse]) error
   	mustEmbedUnimplementedStreamServerServer()
   }
   ```

2. 启动 gRPC 服务器

   ```go
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
   ```

完整代码如下：[server.go](https://github.com/LiangNing7/go-example/blob/main/proto/02-server_stream/server/server.go)

运行服务端：

```bash
$ go run server.go
:8000 net.Listing...
```

### 创建 Client 端

1. 创建调用服务端 `Route` 和`ListValue` 方法

   ```go
   // route 调用服务端 Route 方法.
   func route() {
   	// 创建发送结构体.
   	req := pb.SimpleRequest{
   		Data: "grpc",
   	}
   	// 调用我们的服务(Route方法).
   	// 同时传入了一个 context.Context，在有需要时可以让我们改变 RPC 的行为，
   	// 比如 超时/取消一个正在运行的 RPC.
   	res, err := grpcClient.Route(context.Background(), &req)
   	if err != nil {
   		log.Fatalf("Call Route err: %v", err)
   	}
   	// 打印返回值.
   	log.Println(res)
   }
   
   // listValue() 调用服务端的 ListValue 方法.
   func listValue() {
   	// 创建发送结构体.
   	req := pb.SimpleRequest{
   		Data: "stream server grpc ",
   	}
   
   	// 调用我们的服务(ListValue方法)
   	stream, err := grpcClient.ListValue(context.Background(), &req)
   	if err != nil {
   		log.Fatalf("Call ListStr err: %v", err)
   	}
   	for {
   		// Recv() 方法接收服务端消息，默认每次 Recv() 最大消息长度为 `1024*1024*4`bytes(4M)
   		res, err := stream.Recv()
   		// 判断消息流是否已经结束.
   		if err == io.EOF {
   			break
   		}
   		if err != nil {
   			log.Fatalf("ListStr get stream err: %v", err)
   		}
   		// 打印返回值.
   		log.Println(res.StreamValue)
   	}
   }
   ```

2. 启动 gRPC 客户端

   ```go
   // Address 连接地址.
   const Address string = ":8000"
   
   var grpcClient pb.StreamServerClient
   
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
   	grpcClient = pb.NewStreamServerClient(conn)
   	route()
   	listValue()
   }
   ```

完整代码为：[client.go](https://github.com/LiangNing7/go-example/blob/main/proto/02-server_stream/client/client.go)

运行客户端：

```bash
$ go run client.go
stream server grpc 0
stream server grpc 1
stream server grpc 2
stream server grpc 3
stream server grpc 4
```

客户端不断从服务端获取数据。

### 客户端能自己停止获取数据

1. 先修改服务端的 `ListValue` 方法：

   ```go
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
   ```

2. 再修改客户端的 `listValue` 函数：

   ```go
   // listValue 调用服务端的 ListValue 方法.
   func listValue() {
   	// 创建发送结构体.
   	req := pb.SimpleRequest{
   		Data: "stream server grpc ",
   	}
   
   	// 调用我们的服务 (Route 方法).
   	// 同时传入一个 context.Context，在有需要时可以让我们改变 RPC 的行为，比如超时/取消一个正在运行的 RPC.
   	stream, err := grpcClient.ListValue(context.Background(), &req)
   	if err != nil {
   		log.Fatalf("Call ListStr err: %v", err)
   	}
   	for range 5 {
   		// Recv() 方法接收服务端消息，默认每次 Recv() 最大消息长度为 `1024*1024*4`bytes(4M)
   		res, err := stream.Recv()
   		// 判断消息流是否已经结束.
   		if err == io.EOF {
   			break
   		}
   		if err != nil {
   			log.Fatalf("ListStr get stream err: %v", err)
   		}
   		// 打印返回值.
   		log.Println(res.StreamValue)
   	}
   	// 可以使用 CloseSend() 关闭 stream，这样服务端就不会继续产生流消息.
   	// 调用 CloseSend() 后，若继续调用 Recv()，就会重新激活 stream，接着之前的结果继续获取消息.
   	fmt.Println("暂停调用")
   	stream.CloseSend()
   	fmt.Println("继续调用")
   	for {
   		res, err := stream.Recv()
   		if err == io.EOF {
   			break
   		}
   		if err != nil {
   			log.Fatalf("ListStr get stream err: %v", err)
   		}
   		log.Println(res.StreamValue)
   	}
   }
   ```

只需要调用 `CloseSend()` 方法，就可以关闭服务端的 stream，让它停止发送数据。值得注意的是，调用 `CloseSend()` 后，若继续调用 `Recv()` ，会重新激活stream，接着当前的结果继续获取消息。

这能完美解决客户端`暂停`->`继续`获取数据的操作。

## 客户端流式 RPC

`客户端流式RPC`：与`服务端流式RPC`相反，客户端不断的向服务端发送数据流，而在发送结束后，由服务端返回一个响应。

**情景模拟：客户端大量数据上传到服务端。**

### 新建 proto 文件

完整代码为：[client_stream.proto](https://github.com/LiangNing7/go-example/blob/main/proto/03-client_stream/proto/client_stream.proto)

```protobuf
syntax = "proto3";

package proto;

option go_package = ".;proto";

// 定义发送请求信息
message SimpleRequest{
    // 定义发送的参数，采用驼峰命名方式，小写加下划线，如：student_name
    // 参数类型 参数名 标识号(不可重复)
    string data = 1;
}

// 定义响应信息
message SimpleResponse{
    // 定义接收的参数
    // 参数类型 参数名 标识号(不可重复)
    int32 code = 1;
    string value = 2;
}

// 定义流式请求信息
message StreamRequest{
    //流式请求参数
    string stream_data = 1;
}

// 定义我们的服务（可定义多个服务,每个服务可定义多个接口）
service StreamClient{
    rpc Route (SimpleRequest) returns (SimpleResponse){};

    // 客户端流式rpc，在请求的参数前添加stream
    rpc RouteList (stream StreamRequest) returns (SimpleResponse){};
}
```

最后编译 proto 文件：

进入 `client_stream.proto`文件所在目录，运行：

```bash
$ protoc --go_out=. --go-grpc_out=. ./client_stream.proto
```

### 新建 Server 端

完整代码位于：[server.go](https://github.com/LiangNing7/go-example/blob/main/proto/03-client_stream/server/server.go)

```go
package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "github.com/LiangNing7/go-example/proto/03-client_stream/proto"
	"google.golang.org/grpc"
)

// SimpleService 定义我们的服务.
type SimpleService struct {
	pb.UnimplementedStreamClientServer
}

// Route 实现 Route 方法.
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}

// RouteList 实现 RouteList 方法.
func (s *SimpleService) RouteList(srv pb.StreamClient_RouteListServer) error {
	for {
		// 从流中获取消息.
		res, err := srv.Recv()
		if err == io.EOF {
			// 发送结果，并关闭.
			return srv.SendAndClose(&pb.SimpleResponse{Value: "ok"})
		}
		if err != nil {
			return err
		}
		log.Println(res.StreamData)
	}
}

const (
	// Address 监听地址
	Address string = ":8000"
	// Network 网络通信协议
	Network string = "tcp"
)

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterStreamClientServer(grpcServer, &SimpleService{})

	// 用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
```

### 新建 Client 端

完整代码位于：[client.go](https://github.com/LiangNing7/go-example/blob/main/proto/03-client_stream/client/client.go)

```go
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
```

对应服务端收到的数据为：

```bash
$ go run server.go 
2025/04/24 14:25:04 :8000 net.Listing...
2025/04/24 14:25:06 stream client rpc 0
2025/04/24 14:25:06 stream client rpc 1
2025/04/24 14:25:06 stream client rpc 2
2025/04/24 14:25:06 stream client rpc 3
2025/04/24 14:25:06 stream client\ rpc 4
```

运行客户端：

```bash
$ go run client.go 
2025/04/24 14:25:06 code:200 value:"hello grpc"
2025/04/24 14:25:06 StreamData: 0
2025/04/24 14:25:06 StreamData: 1
2025/04/24 14:25:06 StreamData: 2
2025/04/24 14:25:06 StreamData: 3
2025/04/24 14:25:06 StreamData: 4
2025/04/24 14:25:06 value:"ok"
```

## 双向流式 RPC

### 新建 proto 文件

完整代码位于：[both_stream.proto](https://github.com/LiangNing7/go-example/blob/main/proto/04-both_stream/proto/both_stream.proto)

```protobuf
syntax = "proto3";// 协议为proto3

package proto;

option go_package = ".;proto";

// 定义发送请求信息
message SimpleRequest{
    // 定义发送的参数，采用驼峰命名方式，小写加下划线，如：student_name
    // 参数类型 参数名 标识号(不可重复)
    string data = 1;
}

// 定义响应信息
message SimpleResponse{
    // 定义接收的参数
    // 参数类型 参数名 标识号(不可重复)
    int32 code = 1;
    string value = 2;
}

// 定义流式请求信息
message StreamRequest{
    //流请求参数
    string question = 1;
}

// 定义流式响应信息
message StreamResponse{
    //流响应数据
    string answer = 1;
}

// 定义我们的服务（可定义多个服务,每个服务可定义多个接口）
service Stream{
    rpc Route (SimpleRequest) returns (SimpleResponse){};

    // 双向流式rpc，同时在请求参数前和响应参数前加上stream
    rpc Conversations(stream StreamRequest) returns(stream StreamResponse){};
}
```

最后编译 proto 文件：

进入 `client_stream.proto`文件所在目录，运行：

```bash
$ protoc --go_out=. --go-grpc_out=. ./both_stream.proto
```

### 新建 Server 端

完整代码位于：[server.go](https://github.com/LiangNing7/go-example/blob/main/proto/04-both_stream/server/server.go)

```go
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
```

### 新建 Client 端

完整代码位于：[client.go](https://github.com/LiangNing7/go-example/blob/main/proto/04-both_stream/client/client.go)

```go
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
```

运行服务端：

```bash
$ go run server.go 
2025/04/24 14:30:42 :8000 net.Listing...
2025/04/24 14:30:45 from stream client question: stream client rpc 0
2025/04/24 14:30:45 from stream client question: stream client rpc 1
2025/04/24 14:30:45 from stream client question: stream client rpc 2
2025/04/24 14:30:45 from stream client question: stream client rpc 3
2025/04/24 14:30:45 from stream client question: stream client rpc 4
```

运行客户端：

```bash
$ go run client.go 
2025/04/24 14:30:45 hello grpc
2025/04/24 14:30:45 from stream server answer: the 1 question is stream client rpc 0
2025/04/24 14:30:45 from stream server answer: the 2 question is stream client rpc 1
2025/04/24 14:30:45 from stream server answer: the 3 question is stream client rpc 2
2025/04/24 14:30:45 from stream server answer: the 4 question is stream client rpc 3
2025/04/24 14:30:45 from stream server answer: the 5 question is stream client rpc 4
```

# proto 文件介绍

上面我们大致学习了 gRPC 的使用，现在让我们深入了解 proto 文件。

## 常见关键字

### 基础结构关键字

`syntax`，指定使用的 Protocol Buffers 语法版本，必须放在文件开头。示例如下：

```protobuf
syntax = "proto3";
```

`package`，定义包的命令空间，防止命名冲突。生成的代码会根据语言（如Java、Go）自动映射到对应包结构。示例：

```protobuf
package proto;
```

`import`，导入其他 `.proto` 文件，复用其定义的消息或服务。示例：

```protobuf
import "google/protobuf/empty.proto";
```

### 消息与字段关键字

`message`，定义数据结构，包含一组字段，类似于结构体。示例：

```protobuf
// 定义发送请求信息
message SimpleRequest{
    // 定义发送的参数，采用驼峰命名方式，小写加下划线，如：student_name
    // 参数类型 参数名 标识号(不可重复)
    string data = 1;
}

// 定义响应信息
message SimpleResponse{
    // 定义接收的参数
    // 参数类型 参数名 标识号(不可重复)
    int32 code = 1;
    string value = 2;
}
```

`field`，消息中的字段，其声明方式为：`参数类型 参数名 = 标识号（不可重复）`。参考上面 `message` 中定义的字段。

> tips:
>
> * 每个字段的标签必须唯一，且1-15占用1字节，16-2047占用2字节，通常将常用字段分配1-15。
> * 标识符用于在编译后的二进制消息格式中对字段进行识别，一旦 Protobuf 消息投入使用，字段的标识符就不应再修改。
> * 数字标签的取值范围为 `[1, 536870911]`，其中 19000 至 19999 为保留数字，不能使用。

`repeated`，表示字段是重复的（类似于数组或列表），示例：

```protobuf
repeated string hobbies = 4;
```

`enum`，定义枚举类型。示例：

```protobuf
enum UserRole {
  GUEST = 0;
  USER = 1;
  ADMIN = 2;
}
```

`oneof`，表示一组字段中同一时间只能有一个被设置（类似联合体）。示例：

```protobuf
oneof auth_method {
  string password = 1;
  string token = 2;
}
```

`map`，定义键值对字段。示例：

```protobuf
map<string, int32> scores = 5;
```

`reserved`，保留字段标签或名称，防止后续被误用。示例：

```protobuf
reserved 6 to 10, "field_name";
```

### 服务于方法关键字

`service`，定义 RPC 服务，包含一组方法：

```protobuf
// 定义我们的服务（可定义多个服务，每个服务课定义多个接口）.
service Simple{
  rpc Route (SimpleRequest) returns (SimpleResponse){}
}
```

`rpc`，定义服务中的方法，嘘指定请求和响应类型。[参考上面 RPC 的四种模式](#gRPC 基础)。

### 选项与高级功能

`option`，设置文件、消息、字段或服务的额外配置。常见选项：

* 文件级别：

  ```protobuf
  option go_package = "github.com/LiangNing7/miniblog/pkg/api/apiserver/v1;v1";
  ```

  * 这里我们的 go 项目模块路径为：`github.com/LiangNing7/miniblog`
  * 分号前的部分：`github.com/LiangNing7/miniblog/pkg/api/apiserver/v1`
    这是生成的 Go 代码的**导入路径**（即其他 Go 代码通过该路径导入生成的包）。
  * 分号后的部分：`v1`
    这是生成的 Go 代码的**包名**（即代码中 `package v1` 的声明）。
  * 最后生成的代码被放置在项目的 `pkg/api/apiserver/v1` 目录下。

* 字段级别：

  ```protobuf
  string deprecated_field = 6 [deprecated = true]; // 标记字段为弃用
  ```

## 常见字段定义

### 基础数据类型

| Proto Type |  Go Type  |
| :--------: | :-------: |
|  `double`  | `float64` |
|  `float`   | `float32` |
|  `int32`   |  `int32`  |
|  `int64`   |  `int64`  |
|  `uint32`  | `uint32`  |
|  `uint64`  | `uint64`  |
|  `sint32`  |  `int32`  |
|  `sint64`  |  `int64`  |
| `fixed32`  | `uint32`  |
| `fixed64`  | `uint64`  |
| `sfixed32` |  `int32`  |
| `sfixed64` |  `int64`  |
|   `bool`   |  `bool`   |
|  `string`  | `string`  |
|  `bytes`   | `[]byte`  |
|   `map`    |   `map`   |

在 proto 文件中定义这些数据类型的实例如下：

```protobuf
syntax = "proto3";

package proto;

// 为了生成 Go 代码时使用合适的包名
option go_package = ".;proto";

// 一个包含所有基本字段类型的消息
// 声明格式：参数类型 参数名 = 标识号（不可重复）
message BasicTypes {
  double   field_double    = 1;  // 对应 Go 的 float64
  float    field_float     = 2;  // 对应 Go 的 float32

  int32    field_int32     = 3;  // 对应 Go 的 int32
  int64    field_int64     = 4;  // 对应 Go 的 int64
  uint32   field_uint32    = 5;  // 对应 Go 的 uint32
  uint64   field_uint64    = 6;  // 对应 Go 的 uint64

  sint32   field_sint32    = 7;  // 对应 Go 的 int32（zigzag 编码，适合负数）
  sint64   field_sint64    = 8;  // 对应 Go 的 int64

  fixed32  field_fixed32   = 9;  // 对应 Go 的 uint32（固定 4 字节）
  fixed64  field_fixed64   = 10; // 对应 Go 的 uint64（固定 8 字节）
  sfixed32 field_sfixed32  = 11; // 对应 Go 的 int32（固定 4 字节）
  sfixed64 field_sfixed64  = 12; // 对应 Go 的 int64（固定 8 字节）

  bool     field_bool      = 13; // 对应 Go 的 bool
  string   field_string    = 14; // 对应 Go 的 string
  bytes    field_bytes     = 15; // 对应 Go 的 []byte
  // map<string, int32> 对应 Go 中的 map[string]int32
  map<string, int32> scores = 3;

  // map<string, Status> 对应 Go 中的 map[string]Status
  map<string, Status> status_map = 4;
}
```

### 枚举类型和复合类型

> tips：
>
> proto3 明确规定：
>
> * 枚举必须包含一个值为 `0` 的成员。
> * `0` 值通常作为第一个元素，并建议命名为 `UNSPECIFIED` 或类似名称，以明确表示“未指定”的语义。

```protobuf
syntax = "proto3";

package proto;

// 为了生成 Go 代码时使用合适的包名
option go_package = ".;proto";

// 定义一个枚举类型，用于表示状态
enum Status {
  // 默认的未指定值，Protobuf 要求枚举第一个值必须为 0
  STATUS_UNSPECIFIED = 0;
  STATUS_STARTED     = 1;
  STATUS_IN_PROGRESS = 2;
  STATUS_DONE        = 3;
}

// 演示枚举和 map 的消息
message EnumMapExample {
  // 枚举字段，对应 Go 中的 Status (底层为 int32)
  Status status = 1;

  // 普通 repeated 枚举列表，对应 Go 中的 []Status
  repeated Status all_statuses = 2;
}
```



复合类型：

```protobuf
syntax = "proto3";

package proto;

// 为了生成 Go 代码时使用合适的包名
option go_package = ".;proto";

// Location 消息：表示地理位置坐标
message Location {
  double latitude = 1;   // 纬度（-90 到 90 度，使用 double 保证高精度）
  double longitude = 2;  // 经度（-180 到 180 度）
}

// 行程状态枚举（必须包含 0 值作为默认状态）
enum TripStatus {
  TS_NOT_SPECIFIED = 0;  // 默认/未明确指定的状态（proto3 强制要求）
  NOT_STARTED = 1;       // 行程未开始
  IN_PROGRESS = 2;       // 行程进行中（注意：实际拼写应为 IN_PROGRESS）
  FINISHED = 3;          // 行程已结束
  PAID = 4;              // 行程费用已支付
}

// Trip 消息：表示一个完整的行程记录
message Trip {
  string start = 1;           // 行程起点名称（如 "北京首都机场"）
  string end = 2;             // 行程终点名称（如 "上海浦东国际机场"）
  int64 duration_sec = 3;     // 行程持续时间（秒级精度，避免浮点误差）
  int64 fee_cent = 4;         // 行程费用（以分为单位，避免浮点数精度问题）
  
  Location start_pos = 5;     // 精确的起点坐标（关联 Location 消息）
  Location end_pos = 6;       // 精确的终点坐标
  
  // 行程路径中的多个位置点（repeated 表示数组/列表）
  repeated Location path_locations = 7;

  TripStatus status = 8;      // 行程当前状态（关联 TripStatus 枚举）
}
```

### 自定义类型

ProtoBuf 允许定义自定义类型，例如 **Decimal** 类型。虽然 ProtoBuf 不直接支持 **Decimal** 类型，但可以通过自定义消息类型来实现。

先自定义一个 `DecimalValue` 类型：

```protobuf
syntax = "proto3";

package customTypes;

// 为了生成 Go 代码时使用合适的包名
option go_package = ".;customTypes";

message DecimalValue {
    int32 units = 1;
    int32 nanos = 2;
}
```

然后在其他消息中引用这个自定义类型：

```protobuf
syntax = "proto3";

package proto;

// 为了生成 Go 代码时使用合适的包名
option go_package = ".;proto";

import "customTypes.proto";

message CustomerRebateInfo {
    string str = 1;
    customTypes.DecimalValue total1 = 2;
    customTypes.DecimalValue jinE1 = 3;
}
```

### 二维结构

可以借助使用 `google/protobuf/struct.proto` 文件中的 `Value` 类型，或者使用 Value 包含的 （数值、字符串、布尔、Struct 或 ListValue）类型。

当我们在 Linux 中，使用源码安装 protoc 时，若执行了 `make install`，则其安装目录下会包含一个 `include` 目录，其中就有 `google/protobuf/struct.proto` 文件。

* 在 macOS/Homebrew 安装时，典型路径为：
   `/usr/local/include/google/protobuf/struct.proto`
* 在 Linux 从源码安装时，若执行了 `make install`，则一般在：
   `/usr/local/include/google/protobuf/struct.proto`
   或你的系统头文件目录下的相应位置。

示例如下：[二维结构](https://github.com/LiangNing7/go-example/tree/main/proto/05-two-dimension)

```protobuf
syntax = "proto3";
import "google/protobuf/struct.proto";
option go_package = ".;personpb";
message Foo {
    repeated google.protobuf.Value array = 1;
}
```

在 Go 中使用：

```go
package main

import (
    "encoding/json"
    "fmt"
    structpb "google.golang.org/protobuf/types/known/structpb"
)

type Foo struct {
    Array []*structpb.Value `protobuf:"bytes,1,rep,name=array,proto3" json:"array,omitempty"`
}

func main() {
    l1, _ := structpb.NewList([]any{"1", "2"})
    l2, _ := structpb.NewList([]any{"3", "4"})
    p := Foo{
        Array: []*structpb.Value{
            structpb.NewListValue(l1),
            structpb.NewListValue(l2),
        },
    }

    d, _ := json.Marshal(p)
    fmt.Println(string(d))
}
```

## 编译 proto 文件

一般用 Makefile 来编译 proto 文件：

```makefile
# Protobuf 文件存放路径
APIROOT=$(PROJ_ROOT_DIR)/pkg/api

protoc: # 编译 protobuf 文件.
    @echo "===========> Generate protobuf files"
    @protoc                                              \
        --proto_path=$(APIROOT)                          \
        --proto_path=$(PROJ_ROOT_DIR)/third_party/protobuf    \
        --go_out=paths=source_relative:$(APIROOT)        \
        --go-grpc_out=paths=source_relative:$(APIROOT)   \
        $(shell find $(APIROOT) -name *.proto)
```

`APIROOT`：指向项目中存放 `.proto` 源文件的根目录，用于后续 `protoc` 的输入和输出路径。

`protoc` 命令参数的说明：

* `--proto_path` 或 `-I`：用于指定编译源码的搜索路径，类似于 C/C++中的头文件搜索路径，在构建 `.proto` 文件时，protoc 会在这些路径下查找所需的 Protobuf 文件及其依赖；

  * ```makefile
    --proto_path=$(APIROOT) # 指定项目自己定义的 API 源码目录
    ```

  * ```makefile
    --proto_path=$(PROJ_ROOT_DIR)/third_party/protobuf # 指定第三方依赖路径
    ```

* `--go_out`：用于生成与 gRPC 服务相关的 Go 代码，并配置生成文件的路径和文件结构。例如 `--go_out=plugins=grpc,paths=import:.`。主要参数包括 plugins 和 paths。

  * `plugins=grpc` 表示生成 Go 代码所使用的插件
  * `path=import` 表示生成的 Go 代码的位置。`path` 它支持以下两个选项：
    * `import`（默认值）：按照生成的 Go 代码包的全路径创建目录结构；
    * `source_relative`：表示生成的文件应保持与输入文件相对路径一致。假设 Protobuf 文件位于 `pkg/api/apiserver/v1/example.proto`，启用该选项后，生成的代码也会位于 `pkg/api/apiserver/v1/`目录。如果没有设置 `paths=source_relative`，默认情况下，生成的 Go 文件的路径可能与包含路径有直接关系，并不总是与输入文件相对路径保持一致。

* `--go-grpc_out`：功能与 `--go_out` 类似，但该参数用于指定生成的 `*_grpc.pb.go` 文件的存放路径。

`$(shell find $(APIROOT) -name *.proto)` ，在 Makefile 中通过 `find` 命令递归查找所有后缀为 `.proto` 的文件，作为 `protoc` 的输入列表，这样可以自动拾取新增或删除的 `.proto`，无需手动维护文件列表。

# gRPC 进阶

## 超时设置

gRPC 默认的请求超时时间[默认无限期等待](https://grpc.io/docs/guides/deadlines/#deadlines-on-the-client)，当你没有设置请求超时时间时，所有在运行的请求都占用大量资源且可能运行很长的时间，导致服务资源损耗过高，使得后来的请求响应过慢，甚至会引起整个进程崩溃。

为了避免这种情况，我们的服务应该设置超时时间。前面提到过，当客户端发起请求时候，需要传入上下文 `context.Context`，用于结束`超时`或`取消`的请求。

### proto文件

proto 文件采用 simple.proto，不做变动，代码位于：[simple.proto](https://github.com/LiangNing7/go-example/tree/main/proto/06-deadline/proto)

最后编译 proto 文件：

进入 `simple.proto`文件所在目录，运行：

```bash
$ protoc --go_out=. --go-grpc_out=. ./simple.proto
```

### 创建 Client 端

修改调用服务端方法

1. 把超时时间设置为当前时间 +3 秒

   ```go
   clientDeadline := time.Now().Add(time.Duration(deadlines * time.Second))
   ctx, cancel := context.WithDeadline(ctx, clientDeadline)
   defer cancel()
   ```

2. 响应错误检测中添加超时检测

   ```go
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
   ```

完整的客户端代码位于：[client.go](https://github.com/LiangNing7/go-example/blob/main/proto/06-deadline/client/client.go)

### 创建 Server 端

当请求超时后，服务端应该停止正在进行的操作，避免资源浪费。

```go
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
```

一般地，在写库前进行超时检测，发现超时就停止工作。

完整代码位于：[server.go](https://github.com/LiangNing7/go-example/blob/main/proto/06-deadline/server/server.go)



运行服务端代码：

```bash
$ go run server.go
2025/04/24 19:42:37 :8000 net.Listing...
2025/04/24 19:42:44 context canceled
```

运行客户端代码：

```bash
$ go run client.go
2025/04/24 19:42:44 Route timeout!
```



总结：超时时间的长短需要根据自身服务而定，例如返回一个`hello grpc`，可能只需要几十毫秒，然而处理大量数据的同步操作则可能要很长时间。需要考虑多方面因素来决定这个超时时间，例如系统间端到端的延时，哪些RPC是串行的，哪些是可以并行的等等。

## 认证 - 安全传输

gRPC 是一个典型的 C/S 模型，需要开发客户端和服务端，客户端与服务端需要达成协议，使用某一个确认的传输协议来传输数据，gRPC通常默认是使用`protobuf`来作为传输协议，当然也是可以使用其他自定义的。

那么，客户端与服务端要通信之前，客户端如何知道自己的数据是发给哪一个明确的服务端呢？反过来，服务端是不是也需要有一种方式来弄清楚自己的数据要返回给谁呢？

那么就不得不提 gRPC 的认证

此处说到的认证，不是用户的身份认证，而是指多个`server`和多个`client`之间，如何识别对方是谁，并且可以安全的进行数据传输

- SSL/TLS 认证方式（采用 `http2` 协议）
- 基于 `Token` 的认证方式（基于安全连接）
- 不采用任何措施的连接，这是不安全的连接（默认采用 `http1` )
- 自定义的身份认证

客户端和服务端之间调用，我们可以通过加入证书的方式，实现调用的安全性

TLS (`Transport Layer Security`,安全传输层)，TLS 是建立在传输层 TCP 协议之上的协议，服务于应用层，它的前身是 SSL (`Secure Socket Layer`,安全套接字层)，它实现了将应用层的报文进行加密后再交由TCP进行传输的功能。

TLS 协议主要解决如下三个网络安全问题。

1. 保密(`message privacy`)，保密通过加密 `encryption` 实现，所有信息都加密传输，第三方无法嗅探：
2. 完整性(`message integrity`)，通过 MAC 校验机制，一旦被篡改，通信双方会立刻发现：
3. 认证(`mutual authentication`)，双方认证，双方都可以配备证书，防止身份被冒充

`生产环境可以购买证书或者使用一些平台发放的免费证书`

`key`：服务器上的私钥文件，用于对发送给客户端数据的加密，以及对从客户端接收到数据的解密。
`csr`：证书签名请求文件，用于提交给证书颁发机构(CA)对证书签名。
`crt`：由证书颁发机构(CA)签名后的证书，或者是开发者自签名的证书，包含证书持有人的信息，持有人的公钥，以及签署者的签名等信息。
`pem`：是基于 Base64 编码的证书格式，扩展名包括 PEM、CRT 和 CER。

### 生成 CA 的私钥和证书

创建 `key/` 目录，用于存储密钥和证书，位于：[key/](https://github.com/LiangNing7/go-example/tree/main/proto/07-grpc_security/key)

**生成私钥**

生成 RSA 私钥：

```bash
openssl genrsa -out server.key 2048
```

> 生成 RSA 私钥，命令的最后一个参数，将指定生成密钥的位数，如果没有指定，默认 512

或者生成 ECC 私钥：`openssl ecparam -genkey -name secp384r1 -out server.key`

> 生成 ECC 私钥，命令为椭圆曲线密钥参数生成及操作，本文中 ECC 曲线选择的是 secp384r1

**生成公钥**

```bash
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
```

> openssl req：生成自签名证书，-new指生成证书请求、-sha256指使用sha256加密、-key指定私钥文件、-x509指输出证书、-days 3650为有效期

此后则输入证书拥有者信息：

```bash
#国家名称
Country Name (2 letter code)[AU]:CN
#省名称
State or Province Name (full name)[Some-State]:GuangDong
#城市名称
Locality Name (eg,city)[]:Meizhou
#公司组织名称
Organization Name (eg,company) [Internet widgits Pty Ltd]:Xuexiangban
#部门名称
organizational Unit Name (eg,section) []:go
#服务器0r网站名称
Common Name (e.g.server FQDN or YOUR name) []:liangning
#邮件
Email Address[]：1075090027@qq.com
```

**生成证书签名请求文件**

```bash
openssl req -new -key server.key -out server.csr
```

### 修改 OpenSSL 配置

修改 `openssl.cnf` 以支持SAN（Subject Alternative Name），允许生成包含多域名或通配符的证书。

```bash
# 更改openssl.cnf(windows是openssl.cfg)
#1) 复制一份你安装的openss1的bin目录里面的openssl.cnf文件到你项目所在的目录
#2) 找到[ CA_default ],打开 copy_extensions=copy (就是把前面的#去掉)
#3) 找到[ req ],打开 req_extensions=v3_req # The extensions to add to a certificate request
#4) 找到[ v3_req ],添加 
subjectAltName = @alt_names
#5) 在下面添加新的标签[ alt_names ],和标签字段
[ alt_names ]
DNS.1 = *.liangning.com
```

![image-20250424211425440](http://images.liangning7.cn/typora/202504242114893.png)

### 生成终端实体的私钥和证书

```bash
#生成证书私钥test.key
openssl genpkey -algorithm RSA -out test.key

#通过私钥test.key生成证书请求文件test.csr (注意cfg和cnf)
openssl req -new -nodes -key test.key -out test.csr -subj "/C=cn/OU=myorg/O=mycomp/CN=myname" -config ./openssl.cnf -extensions v3_req
#test.csr是上面生成的证书请求文件。ca.crt/server.key是CA证书文件和key,用来对test.csr进行签名认证。这两个文件在第一部分生成。

#生成SAN证书pem
openssl x509 -req -days 365 -in test.csr -out test.pem -CA server.pem -CAkey server.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req
```

创建完成后，`key`文件夹下有如下8个文件：

```bash
$ ls
openssl.cnf  server.csr  server.key  server.pem  server.srl  test.csr  test.key  test.pem
```

### 使用简单模式 RPC 的 proto 文件

还是使用简单模式 RPC 的 proto 文件：[simple.proto](https://github.com/LiangNing7/go-example/blob/main/proto/07-grpc_security/proto/simple.proto)

最后编译 proto 文件：

进入 `simple.proto`文件所在目录，运行：

```bash
$ protoc --go_out=. --go-grpc_out=. ./simple.proto
```

### 创建 Server 端

```go
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
	// 新建 gRPC 实例，并开启 TLS 认证.
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// 在 gRPC 服务器中注册我们的服务.
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	log.Println(Address + " net.Listing with TLS and token...")
	// 用服务器 Server() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用.
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
```

* `credentials.NewServerTLSFromFile`：从输入证书文件和密钥文件为服务端构造TLS凭证
* `grpc.Creds`：返回一个ServerOption，用于设置服务器连接的凭证。

完整代码位于：[server.go](https://github.com/LiangNing7/go-example/blob/main/proto/07-grpc_security/server/server.go)

### 创建 Client 端

```go
var grpcClient pb.SimpleClient

func main() {
	// 从输入的证书文件中为客户端构造 TLS 凭证.
	creds, err := credentials.NewClientTLSFromFile("../key/test.pem", "blog.liangning7.cn")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	// 连接服务器.
	conn, err := grpc.NewClient(
		Address,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立 gRPC 连接.
	grpcClient = pb.NewSimpleClient(conn)
	route()
}
```

* `credentials.NewClientTLSFromFile`：从输入的证书文件中为客户端构造TLS凭证。
* `grpc.WithTransportCredentials`：配置连接级别的安全凭证（例如，TLS/SSL），返回一个 DialOption，用于连接服务器。

完整代码位于：[client.go](https://github.com/LiangNing7/go-example/blob/main/proto/07-grpc_security/client/client.go)



运行服务端：

```bash
$ go run server.go
2025/04/24 21:43:13 :8000 net.Listing with TLS and token...
```

运行客户端：

```bash
$ go run client.go
2025/04/24 21:44:45 code:200  value:"hello grpc"
```

到这里，已经完成 TLS 证书认证了，gRPC 传输不再是明文传输。此外，添加自定义的验证方法能使 gRPC 相对更安全。

## Token 认证

实现 Token 认证，即使用 gRPC  中定义的自定义认证接口，将所需的安全认证信息添加到每个RPC方法的上下文中。

```go
// PerRPCCredentials defines the common interface for the credentials which need to
// attach security information to every RPC (e.g., oauth2).
type PerRPCCredentials interface {
	// GetRequestMetadata gets the current request metadata, refreshing tokens
	// if required. This should be called by the transport layer on each
	// request, and the data should be populated in headers or other
	// context. If a status code is returned, it will be used as the status for
	// the RPC (restricted to an allowable set of codes as defined by gRFC
	// A54). uri is the URI of the entry point for the request.  When supported
	// by the underlying implementation, ctx can be used for timeout and
	// cancellation. Additionally, RequestInfo data will be available via ctx
	// to this call.  TODO(zhaoq): Define the set of the qualified keys instead
	// of leaving it as an arbitrary string.
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	// RequireTransportSecurity indicates whether the credentials requires
	// transport security.
	RequireTransportSecurity() bool
}
```

* `GetRequestMetadata`：获取元数据，也就是客户端提过的`key:value`对，`context`用于控制超时和取消，`uri`是请求入口出的`uri`
* `RequireTransportSecurity`：是否需要基于 TLS 认证进行安全传输，如果返回值为 `true`，则必须加上 TLS 验证，返回值是 `false `则不用。



gRPC 的 Token 认证（调用凭证，Call Credentials）**不需要必须与 TLS（传输层安全协议，Transport Credentials）一起使用**。但从安全角度出发，**强烈建议始终结合 TLS**。因为通过 `Insecure` 通道（不启用 TLS）发送 Token（例如 JWT、OAuth2 Token）。此时，Token 会以明文形式传输，服务端可以通过拦截器（Interceptor）或中间件解析 Token 进行身份验证。**TLS 的作用**：加密通信链路，防止 Token 和数据被窃听、篡改或伪造。

这里 Token 认证的传输方式，我们基于 TLS 传输，先直接复制 TLS 传输的代码。

### 实现 `PerRPCCredentials` 接口

代码位于：[auth.go](https://github.com/LiangNing7/go-example/blob/main/proto/08-token/auth/auth.go)

```go
package auth

import "context"

// Token token 认证.
type Token struct {
	AppID     string
	AppSecret string
}

// GetRequestMetadata 获取当前请求认证所需的元数据.
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": t.AppID, "app_secret": t.AppSecret}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输.
func (t *Token) RequireTransportSecurity() bool {
	return true
}
```

### 客户端请求添加 Token 到上下文

客户端在调用 `grpc.NewClient()` 时添加自定义验证方法进去：

```go
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
```

完整代码位于：[client.go](https://github.com/LiangNing7/go-example/blob/main/proto/08-token/client/client.go)

### 服务端验证 Token

首先需要从上下文中获取元数据，然后从元数据中解析 Token 进行验证。

```go
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

// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	// 检测 Token 是否有效.
	if err := Check(ctx); err != nil {
		return nil,err
	}
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}
```

* `metadata.FromIncomingContext`：从上下文中获取元数据

完整代码位于：[server.go](https://github.com/LiangNing7/go-example/blob/main/proto/08-token/server/server.go)

### 服务端添加拦截器

服务端代码中，每个服务的方法都需要添加 `Check(ctx)` 来验证 Token，这样十分麻烦。gRPC 拦截器，能很好地解决这个问题。

gRPC 拦截器是一个 Web 中间件。利用拦截器，开发者可以在不侵入业务逻辑的前提下修改或者记录服务端或客户端的请求与响应，利用拦截器可以实现诸如日志记录、权限认证、限流等诸多功能。

gRPC 的通信模式分为 Unary 和 Streaming 两种模式，拦截器也分为两种：UnaryInterceptor（一元拦截器）和 StreamInterceptor（流式拦截器）。这两种拦截器可以分别应用在服务端和客户端，所以 gRPC 框架中，一共提供了四种拦截器：

* **UnaryServerInterceptor：**服务端一元拦截器，适用于简单 RPC 调用。它会在服务端接收到请求时执行拦截逻辑，通常用于对请求进行预处理、授权、认证、日志记录、错误处理等；
* **StreamServerInterceptor：**服务端流式拦截器，适用于流式 RPC 调用，例如客户端流式、服务端流式和双向流式 RPC 调用。它会在服务端接收到流式请求时进行拦截，允许开发者对流式数据进行操作和处理；
* **UnaryClientInterceptor：**客户端一元拦截器，适用于简单 RPC 调用（Unary RPC）。它会拦截客户端发起的调用，通常用于操控请求或响应，比如：请求重试、请求参数的统一注入、加密、客户端的日志记录等；
* **StreamClientInterceptor：**客户端流式拦截器，适用于流式 RPC 调用（客户端流式、服务端流式、双向流式）。它允许在流式调用时通过拦截客户端流（ClientStream）创建过程自定义逻辑，开发者可以围绕流式数据进行操作。

这里使用服务端一元拦截器：

```go
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
```

* `grpc.UnaryServerInterceptor`：为一元拦截器，只会拦截简单RPC方法。流式RPC方法需要使用流式拦截器 `grpc.StreamInterceptor` 进行拦截，使用了该拦截器后，在每个服务的方法就不需要添加 `Check(ctx)` 来验证 Token。

启动服务端：

```bash
$ go run server.go
2025/04/24 22:35:51 :8000 net.Listing with TLS and token...
```

启动客户端：

```bash
$ go run client.go
2025/04/24 22:35:57 code:200  value:"hello grpc"
```

客户端发起请求，当 Token 不正确时候，例如，当 `app_id=grpc` 时，会返回

```bash
2025/04/24 22:36:33 Call Route err: rpc error: code = Unauthenticated desc = Token 无效: app_id=grpc, app_secret=123456
```

# gRPC 中间件使用

`go-grpc-middleware` 封装了认证（auth）, 日志（ logging）, 消息（message）, 验证（validation）, 重试（retries） 和监控（retries）等拦截器。

安装：`go get github.com/grpc-ecosystem/go-grpc-middleware`

使用：

```go
import "github.com/grpc-ecosystem/go-grpc-middleware"
myServer := grpc.NewServer(
  grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        grpc_ctxtags.StreamServerInterceptor(),
        grpc_opentracing.StreamServerInterceptor(),
        grpc_prometheus.StreamServerInterceptor,
        grpc_zap.StreamServerInterceptor(zapLogger),
        grpc_auth.StreamServerInterceptor(myAuthFunction),
        grpc_recovery.StreamServerInterceptor(),
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_ctxtags.UnaryServerInterceptor(),
        grpc_opentracing.UnaryServerInterceptor(),
        grpc_prometheus.UnaryServerInterceptor,
        grpc_zap.UnaryServerInterceptor(zapLogger),
        grpc_auth.UnaryServerInterceptor(myAuthFunction),
        grpc_recovery.UnaryServerInterceptor(),
    )),
)
```

* `grpc.StreamInterceptor`：添加流式RPC的拦截器。
* `grpc.UnaryInterceptor`：添加简单RPC的拦截器。

## grpc_zap 日志记录

在 `server/` 目录下创建 `middleware/` 目录，然后创建 `zap` 包，位于：[zap.go](https://github.com/LiangNing7/go-example/tree/main/proto/09-middleware/server/middleware/zap)

先创建 `zap.Logger` 实例，

```go
// ZapInterceptor 返回 zap.logger 实例.
func ZapInterceptor() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  "log/debug.log",
		MaxSize:   1024, // MB
		LocalTime: true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.NewAtomicLevel(),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
}
```

把 zap 拦截器添加到服务端

```go
grpcServer := grpc.NewServer(
	cred.TLSInterceptor(),
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
	)),
	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
	)),
)
```

日志的各个字段如下：

```log
{
	  "level": "info",						// string  zap log levels
	  "msg": "finished unary call",					// string  log message

	  "grpc.code": "OK",						// string  grpc status code
	  "grpc.method": "Ping",					/ string  method name
	  "grpc.service": "mwitkow.testproto.TestService",              // string  full name of the called service
	  "grpc.start_time": "2006-01-02T15:04:05Z07:00",               // string  RFC3339 representation of the start time
	  "grpc.request.deadline": "2006-01-02T15:04:05Z07:00",         // string  RFC3339 deadline of the current request if supplied
	  "grpc.request.value": "something",				// string  value on the request
	  "grpc.time_ms": 1.345,					// float32 run time of the call in ms

	  "peer.address": {
	    "IP": "127.0.0.1",						// string  IP address of calling party
	    "Port": 60216,						// int     port call is coming in on
	    "Zone": ""							// string  peer zone for caller
	  },
	  "span.kind": "server",					// string  client | server
	  "system": "grpc",						// string

	  "custom_field": "custom_value",				// string  user defined field
	  "custom_tags.int": 1337,					// int     user defined tag on the ctx
	  "custom_tags.string": "something"				// string  user defined tag on the ctx
}
```

## grpc_auth 认证

go-grpc-middleware 中的 grpc_auth 默认使用 `authorization` 认证方式，以 authorization 为头部，包括 `basic`，`bearer` 形式等。下面介绍 `bearer token` 认证。`bearer` 允许使用 `access key`（如 JWT）进行访问。

新建 grpc_auth 服务端拦截器，位于：[auth.go](https://github.com/LiangNing7/go-example/blob/main/proto/09-middleware/server/middleware/auth/auth.go)

```go
package auth

import (
	"context"
	"errors"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Token 用户信息.
type TokenInfo struct {
	ID    string
	Roles []string
}

// AuthInterceptor 认证拦截器，对以 authorization 为头部，
// 形式为 `bearer token` 的 Token 进行验证.
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, " %v", err)
	}
	// 使用 context.WithValue 添加了值后，可以使用 Value(key) 方法获取值.
	newCtx := context.WithValue(ctx, tokenInfo.ID, tokenInfo)
	return newCtx, nil
}

// 解析 token，并进行验证.
func parseToken(token string) (TokenInfo, error) {
	var tokenInfo TokenInfo
	if token == "grpc.auth.token" {
		tokenInfo.ID = "1"
		tokenInfo.Roles = []string{"admin"}
		return tokenInfo, nil
	}
	return tokenInfo, errors.New("Token 无效: bearer " + token)
}

// 从 token 中获取用户唯一标识.
func userClaimsFromToken(tokenInfo TokenInfo) string {
	return tokenInfo.ID
}
```

代码中的对 token 进行简单验证并返回模拟数据。

客户端请求添加 `bearer token`

gRPC 中默认定义了 `PerRPCCredentials`，是提供用于自定义认证的接口，它的作用是将所需的安全认证信息添加到每个 RPC 方法的上下文中。其包含 2 个方法：

* `GetRequestMetadata`：获取当前请求认证所需的元数据
* `RequireTransportSecurity`：是否需要基于 TLS 认证进行安全传输

接下来我们实现这两个方法，完整代码位于：[auth.go](https://github.com/LiangNing7/go-example/blob/main/proto/09-middleware/client/auth/auth.go)

```go
package auth

import "context"

// Token token认证
type Token struct {
	Value string
}

const headerAuthorize string = "authorization"

// GetRequestMetadata 获取当前请求认证所需的元数据
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{headerAuthorize: t.Value}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输
func (t *Token) RequireTransportSecurity() bool {
	return true
}
```

发送请求时添加 token

```go
func main() {
	// 从输入的证书文件中为客户端构造 TLS 凭证.
	creds, err := credentials.NewClientTLSFromFile("../key/test.pem", "blog.liangning7.cn")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	// 构建 Token.
	token := auth.Token{
		Value: "bearer grpc.auth.token",
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
```

把 grpc_auth 拦截器添加到服务端

```go
grpcServer := grpc.NewServer(cred.TLSInterceptor(),
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	        grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		    grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
		)),
	)
```

写到这里，服务端都会拦截请求并进行`bearer token`验证，使用`bearer token`是规范了与`HTTP`请求的对接，毕竟gRPC也可以同时支持`HTTP`请求。

## grpc_recovery 恢复

把gRPC中的`panic`转成`error`，从而恢复程序。

自定义错误返回：[recovery.go](https://github.com/LiangNing7/go-example/blob/main/proto/09-middleware/server/middleware/recovery/recovery.go)

```go
package recovery

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryInterceptor panic 时返回 Unknown 错误码.
func RecoveryInterceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p any) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	})
}
```

添加 grpc_recovery 拦截器到服务端

```go
grpcServer := grpc.NewServer(cred.TLSInterceptor(),
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	        grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
			grpc_recovery.StreamServerInterceptor(recovery.RecoveryInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		    grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
            grpc_recovery.UnaryServerInterceptor(recovery.RecoveryInterceptor()),
		)),
	)
```

## 总结

代码详细请移步：[代码](https://github.com/LiangNing7/go-example/tree/main/proto/09-middleware)

这里介绍了 `go-grpc-middleware` 中的 `grpc_zap`、`grpc_auth` 和 `grpc_recovery` 拦截器的使用。`go-grpc-middleware` 中其他拦截器可参考[GitHub](https://github.com/grpc-ecosystem/go-grpc-middleware)学习使用。

# grpc-gateway

grpc-gateway 是 protoc 的一个插件。它读取 gRPC 服务定义，并生成反向代理服务器（Reverse Proxy）。反向代理服务器根据 gRPC 服务定义中的 `google.api.http` 注释生成，能够将 RESTful JSON API 转换为 gRPC 请求，从而实现同时支持 gRPC 客户端和 HTTP 客户端调用 gRPC 服务的功能。下图展示了通过 gRPC 请求和 REST 请求调用 gRPC 服务的流程。

![image-20250423130913933](http://images.liangning7.cn/typora/202504231309086.png)

在传统的 gRPC 应用程序中，通常会创建一个 gRPC 客户端与 gRPC 服务进行交互。但在此场景中，并未直接构建 gRPC 客户端，而是利用 grpc-gateway 构建了一个反向代理服务。该代理服务为 gRPC 服务中的每个远程方法暴露了 RESTful API，并接收来自 REST 客户端的 HTTP 请求。随后，它将 HTTP 请求转换为 gRPC 消息，并调用后端服务的远程方法。后端服务返回的响应消息会被代理服务再次转换为 HTTP 响应，并发送回客户端。

## 为什么需要 gRPC-Gateway

在 Go 项目开发中，为了提升接口性能并便于内部系统之间的接口调用，通常会使用 RPC 协议通信。而对于外部系统，为了提供更标准、更通用且易于理解的接口调用方式，往往会使用与编程语言无关的 HTTP 协议进行通信。这两种不同的协议在代码实现上存在较大差异。如果开发者希望同时实现内部系统使用 RPC 协议通信以及外部系统通过 HTTP 协议访问，则需要维护两套服务及接口实现代码。这将显著增加后期代码维护与升级的成本，同时也容易导致错误。如果能够将 HTTP 请求转换为 gRPC 请求，并统一通过 gRPC 接口实现所有功能，那么上述问题即可迎刃而解。gRPC-Gateway 正是通过类似的实现方式解决了这一问题。

## 如何使用 grpc-gateway

grpc-gateway 是 protoc 工具的一个插件，所以首先需要确保系统已经安装了 protoc 工具。此外，使用 gprc-gateway 还需要安装以下两个插件：

* `protoc-gen-grpc-gateway`：为 gRPC 服务生成 HTTP/REST API 反向代理代码，从而实现对 gRPC 服务的 HTTP 映射支持；
* `protoc-gen-openapiv2`：用于从 Protobuf 描述中生成 OpenAPI v2（Swagger）定义文件。

插件安装命令如下：

```bash
$ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
$ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0
```

由于之前已经实现了一个简单的 gRPC 服务，现在只需要在 MiniBlog 服务定义文件中添加 gRPC-Gateway 注解即可，注解定义了 gRPC 服务如何映射到 RESTful JSON API，包括指定 HTTP 请求方法、请求路径、请求参数等信息。MiniBlog gRPC 服务 UpdatePost 接口的注解如下所示：

```go
// 定义了一个 MiniBlog RPC 服务
service MiniBlog {
    // UpdatePost 更新文章
    rpc UpdatePost(UpdatePostRequest) returns (UpdatePostResponse) {
        // 将 UpdatePost 映射为 HTTP PUT 请求，并通过 URL /v1/posts/{postID} 访问
        // {postID} 是一个路径参数，grpc-gateway 会根据 postID 名称，将其解析并映射到
        // UpdatePostRequest 类型中相应的字段.
        // body: "*" 表示请求体中的所有字段都会映射到 UpdatePostRequest 类型。
        option (google.api.http) = {
            put: "/v1/posts/{postID}",
            body: "*",
        };

        // 提供用于生成 OpenAPI 文档的注解
        option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
            // 在文档中简要描述此操作的功能：更新文章。
            summary: "更新文章";
            // 为此操作指定唯一标识符（UpdatePost），便于跟踪
            operation_id: "UpdatePost";
            // 将此操作归类到 "博客管理" 标签组，方便在 OpenAPI 文档中组织接口分组
            tags: "博客管理";
        };
    }
}

// UpdatePostRequest 表示更新文章请求
message UpdatePostRequest {
    // postID 表示要更新的文章 ID，对应 {postID}
    string postID = 1;
    // title 表示更新后的博客标题
    optional string title = 2;
    // content 表示更新后的博客内容
    optional string content = 3;
}

// UpdatePostResponse 表示更新文章响应
message UpdatePostResponse {
}
```

在 `UpdatePost` 接口定义中，使用 `google.api.http` 注解，将 `UpdatePost` 映射为 `HTTP PUT` 请求，并通过 `URL /v1/posts/{postID}` 访问。`{postID}` 是一个路径参数，`grpc-gateway` 会根据 `postID` 名称，将其解析并映射到 `UpdatePostRequest` 类型中的 `postID` 字段。`body: "*"` 表示请求体中的所有字段都会映射到 `UpdatePostRequest` 类型中的同名字段中。

在通过 `google.api.http` 注解将 gRPC 方法映射为 HTTP 请求时，有以下规则需要遵守：

1. HTTP 路径可以包含一个或多个 gRPC 请求消息中的字段，但这些字段应该是 nonrepeated 的原始类型字段；
2. 如果没有 HTTP 请求体，那么出现在请求消息中但没有出现在 HTTP 路径中的字段，将自动成为 HTTP 查询参数；
3. 映射为 URL 查询参数的字段应该是原始类型、repeated 原始类型或 nonrepeated 消息类型；
4. 对于查询参数的 repeated 字段，参数可以在 URL 中重复，形式为 `…?param=A&m=B`；
5. 对于查询参数中的消息类型，消息的每个字段都会映射为单独的参数，比如 `…?foo.a=A&foo.b=B&foo.c=C`。

此外，还可以根据需要添加全局的 OpenAPI 配置，用于在生成 OpenAPI 文档时，提供更详细的配置信息。

至此，大致完成了 gRPC 的学习！！！
