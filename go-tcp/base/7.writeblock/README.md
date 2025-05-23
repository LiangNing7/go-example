写阻塞和写入部分数据的场景

> “写入部分数据”需要在分别启动服务端和客户端后，当客户端写入完第一批数据阻塞后，杀掉服务端进程，客户端会报 `write: broken pipe`的错误

先启动服务端：`go run base/7.writeblock/client/main.go`

再启动客户端：`go run base/7.writeblock/server/main.go`

**客户端启动后的行为**

当你运行客户端代码时，将会看到如下过程：

* **连接到服务端**
   客户端首先使用 `net.Dial("tcp", ":8888")` 连接到服务端。如果连接成功，程序输出“dial ok”。
* **无限写入数据**
   客户端构造了一个大小为 65536 字节的 `data` 字节片（即大约 64KB），并进入一个无限循环，不断调用 `conn.Write(data)`。
   每次写入成功，客户端会累加写入总字节数并打印“write ... bytes this time, ... bytes in total”。

## 连接建立后的服务端行为

进入 `handleConn` 函数后，会依次发生以下动作：

* **初始延时**
   函数开始时有 `time.Sleep(time.Second * 10)`，这意味着服务端在和客户端建立连接后会先等待 10 秒。这 10 秒内服务端并不会立即读取客户端发送的数据。
* **读取数据前的等待**
   延时结束后，服务端打印“准备读取...”。随后进入一个无限循环，每次循环开始前会额外延时 5 秒。这意味着服务端每 5 秒才会调用一次 `c.Read(buf)` 进行数据读取。
* **数据读取**
   每次调用 `c.Read(buf)` 后，根据收到的数据量，服务端会打印时间戳和本次读取的字节数。如果遇到错误（比如连接断开或其他网络错误），则打印错误信息并终止该 goroutine。