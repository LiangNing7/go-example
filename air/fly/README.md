# 效率工具：使用 Air 热加载 Go 应用程序

在项目开发阶段，热重载技术非常有用，通过热重载，可以实现在无需手动干预的情况下，修改代码文件后，自动重启 Go 应用。这极大的提升了开发体验，同时也节约了我们的开发时间。

## 简介

Air 是为 Go 应用开发设计的一款支持热重载的命令行工具。

以下是 air 官方总结的特色：

* 彩色的日志输出
* 自定义构建或必要的命令
* 支持外部子目录
* 在 Air 启动之后，允许监听新创建的路径
* 更棒的构建过程

## 快速开始

### 安装

首先我们来安装 air，命令如下。

> note:
>
> 建议使用 go1.23 或更高版本。

```bash
$ go install github.com/air-verse/air@latest
```

### 使用

准备以下目录结构。

```bash
$ tree -F fly
fly/
├── README.md
└── main.go
```

现在我们来编写代码体验一下 air，`main.go` 代码如下。

```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 注册优雅退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 初始化 HTTP 服务器
	server := &http.Server{Addr: ":8080"}

	// 定义路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🚀 Hello Air! (PID: %d)", os.Getpid()) // PID 用于验证热替换
	})

	// 启动服务协程
	go func() {
		fmt.Printf("Server started at http://localhost:8080 (PID: %d)\n", os.Getpid())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// 阻塞等待终止信号
	<-signalChan
	fmt.Println("Server shutting down...")
	server.Close()
}
```

这是一个简单的 Web Server 程序，根路由 `/` 的处理函数中获取并返回了当前程序的 `pid`。

在项目根目录 `fly/` 下直接执行 `air` 命令。

![image-20250418170528736](http://images.liangning7.cn/typora/202504181705796.png)

可以看到，只需要简单的 `air` 命令，我们的 Go 程序就启动起来了。

打开新的终端，使用 `curl` 命令访问 http://localhost:8080。

```bash
curl http://localhost:8080
🚀 Hello Air! (PID: 859842)
```

![image-20250418170838419](http://images.liangning7.cn/typora/202504181708480.png)

重新使用 `curl` 命令访问 http://localhost:8080。

```bash
curl http://localhost:8080
🚀 Hello Fly! (PID: 859985)
```

响应结果已经变了，并且程序的 `pid` 也发生了变化，说明 Go 程序真的被重启了。

那么接下来，我们再将如下内容，写入 `README.md` 中并保存文件。

````markdown
# Air 热重载演示项目

## 快速开始

### 1. 安装 Air

```bash
$ go install github.com/air-verse/air@latest
```

### 2. 使用 Air

```bash
$ cd /path/to/your_project
$ air
```
````

可以发现，这次 `air` 并没有重启 Go 程序。

![image-20250418171307598](http://images.liangning7.cn/typora/202504181713664.png)

想来这也合情合理，`README.md` 文件并不是 Go 程序代码，不会影响程序执行结果，所以无需重启。

不过有一种情况则需要考虑，如果我们提供一个接口，可以返回静态的 `README.md` 文件内容，那么 `README.md` 文件修改，就需要重启 Go 程序。

现在，我们按下 `Ctrl + C` 结束进程。

![image-20250418171354374](http://images.liangning7.cn/typora/202504181713438.png)

细心的你也许已经发现，我们的 `main.go` 代码在实现优雅退出后会输出一条日志 `Server shutting down...`，可是现在终端中并没有输出。

要如何解决这个问题呢？咱们接着往下看。

## 使用进阶

### 自定义配置

air 支持很多高级功能，这些功能都可以通过配置文件中的配置项来开启或关闭。

air 的配置文件格式问 `toml`，我们有两种方式来获得 air 配置文件模板。

一种是 air 官方提供了 example 配置样例，你可以在此查看 [Air Example](https://github.com/air-verse/air/blob/master/air_example.toml)。

另外一种是使用 `air init` 命令来自动生成。

```bash
$ air init

  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.61.7, built with Go go1.24.0

.air.toml file created to the current directory with the default settings
```

在项目根目录下执行 `air init` 后，会生成叫 `.air.toml` 的配置文件。

air 的全量配置内容如下，我对其做了详细的中文注释。

```bash
# Air 热重载工具的 TOML 格式配置文件
# 完整文档参考：https://github.com/air-verse/air

# 工作目录
# 支持相对路径（.）或绝对路径，注意后续目录必须在此目录下
root = "."
# air 执行过程中产生的临时文件存储目录
tmp_dir = "tmp"

[build]
# 构建前执行的命令列表（每项命令按顺序执行）
pre_cmd = ["echo 'hello air' > pre_cmd.txt"]
# 主构建命令（支持常规 shell 命令或 make 工具）
cmd = "go build -o ./tmp/main ."
# 构建后执行的命令列表（相当于按 ^C 程序终止后触发）
post_cmd = ["echo 'hello air' > post_cmd.txt"]
# 从 `cmd` 编译生成的二进制文件路径
bin = "tmp/main"
# 自定义运行参数（可设置环境变量）
full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
# 传递给二进制文件的运行参数（示例将执行 './tmp/main hello world'）
args_bin = ["hello", "world"]
# 监听以下扩展名的文件变动
include_ext = ["go", "tpl", "tmpl", "html"]
# 排除监视的目录列表
exclude_dir = ["assets", "tmp", "vendor", "frontend/node_modules"]
# 指定要监视的目录（空数组表示自动检测）
include_dir = []
# 指定要监视的特定文件（空数组表示自动检测）
include_file = []
# 排除监视的特定文件（空数组表示不过滤）
exclude_file = []
# 通过正则表达式排除文件（示例排除所有测试文件）
exclude_regex = ["_test\\.go"]
# 是否排除未修改的文件（提升性能）
exclude_unchanged = true
# 是否跟踪符号链接目录，允许 Air 跟踪符号链接（软链接）指向的目录/文件变化，适用于项目依赖外部符号链接资源的场景
follow_symlink = true
# 日志文件存储路径（位于 tmp_dir 下）
log = "air.log"
# 是否使用轮询机制检测文件变化（替代 fsnotify），Air 默认使用的跨平台文件监控库（基于 Go 的 fsnotify 包），通过操作系统事件实时感知文件变化
poll = false
# 轮询检测间隔（默认最低 500ms）
poll_interval = 500 # ms
# 文件变动后的延迟构建时间（防止高频触发）
delay = 0 # ms
# 构建出错时是否终止旧进程
stop_on_error = true
# 是否发送中断信号再终止进程（Windows 不支持）
send_interrupt = false
# 发送中断信号后的终止延迟
kill_delay = 500 # nanosecond
# 当程序退出时，是否重新运行二进制文件（适合 CLI 工具）
rerun = false
# 重新运行的时间间隔
rerun_delay = 500

[log]
# 是否显示日志时间戳
time = false
# 仅显示主日志（过滤监控/构建/运行日志）
main_only = false
# 禁用所有日志输出
silent = false

[color]
# 主日志颜色（支持 ANSI 颜色代码）
main = "magenta"
# 文件监控日志颜色
watcher = "cyan"
# 构建过程日志颜色
build = "yellow"
# 运行日志颜色
runner = "green"

[misc]
# 退出时自动清理临时目录（tmp_dir）
clean_on_exit = true

[screen]
# 重建时清空控制台界面
clear_on_rebuild = true
# 保留滚动历史（不清屏时有效）
keep_scroll = true

[proxy]
# 启用浏览器实时重载功能
# 参考：https://github.com/air-verse/air/tree/master?tab=readme-ov-file#how-to-reload-the-browser-automatically-on-static-file-changes
enabled = true
# 代理服务器端口（Air 监控端口），浏览器连接到 proxy_port，Air 将请求转发到应用的真实端口 app_port
proxy_port = 8090
# 应用实际运行端口（需与业务代码端口一致）
app_port = 8080
```

我们可以发现配置项是按功能进行分类的，有如下几块配置。

* 全局配置：`root` 和 `tmp_dir` 分别表示项目的工作目录和临时目录。
* `build` 类配置：都是与构建相关的配置项。
* `log` 类配置：与日志相关的配置项。
* `color` 类配置：与输出颜色相关的配置项。
* `misc` 类配置：杂项配置，其实只有一个 `clean_on_exit` 配置可以用来清理临时目录。
* `screen` 类配置：控制台相关的配置项。
* `proxy` 类配置：代理相关的配置项。

其实分析完了这些配置项，你就能够发现，我们常用的配置其实也就是 `build` 配置项下的那几个，其他的等用到了再去研究不迟。

我们通常可以重点关注这几个配置项：

```bash
[build]
# 主构建命令（支持常规 shell 命令或 make 工具）
cmd = "go build -o ./tmp/main ."
# 从 `cmd` 编译生成的二进制文件路径
bin = "tmp/main"
# 自定义运行参数（可设置环境变量）
full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
# 传递给二进制文件的运行参数（示例将执行 './tmp/main hello world'）
args_bin = ["hello", "world"]
```

`cmd`、`bin`、`args_bin` 3 者正是对应了 Go 应用构建、执行、和传递命令行参数的功能。所以这也是最常用的 3 项配置。`full_bin` 则是高阶版的 `bin`，可以设置环境变量。

这里需要注意，`cmd` 不仅支持 `go build` 构建命令，它还支持常规的 `shell` 命令和 `make` 命令，这为应用构建提供了更多的灵活性。

此外 `misc` 配置项下的 `clean_on_exit` 配置也比较有用，将其设置为 `true` 可以自动清理 `tmp_dir` 目录产生的临时文件。

学习了 air 的配置项，接下来我们来解决程序退出日志 `Server shutting down...` 未输出的问题。

首先将以上配置信息保存在项目根目录 `fly/` 下的 `.air.toml` 文件中。

接着，修改如下这两项配置。

```bash
[build]
  kill_delay = "1s"
  send_interrupt = true
```

开启 `send_interrupt` 配置项后，`air` 在终止我们的 Go 程序之前，会向其发送 `Ctrl + C` 信号，`kill_delay` 配置项保证 `air` 命令为 Go 程序预留足够多的退出时间，给程序优雅退出的机会。

再次使用 `air` 命令启动程序。

![image-20250418172929111](http://images.liangning7.cn/typora/202504181729216.png)

按下 `Ctrl + C` 结束进程，同样会得到 `Server shutting down...` 日志输出。并且等待 1s 过后 `air` 才会输出 `see you again~` 并退出。

至于其他配置项，就交给你自行去探索了。

### 命令行参数

`air` 命令不仅支持使用配置文件来开启或关闭功能，它也支持直接通过命令行参数的方式来开启或关闭某项功能。

在前文中我们已经使用过 `air init` 来创建配置文件，现在来看看 `air` 还支持哪些命令行参数。

```bash
$  air -h  
Usage of air:

If no command is provided air will start the runner with the provided flags

Commands:
  init  creates a .air.toml file with default settings to the current directory

Flags:
  -build.args_bin string
        Add additional arguments when running binary (bin/full_bin).
  ...
  -c string
        config path
  ...
  -v    show version
```

执行 `air -h` 或 `air --help` 即可查看命令行帮助信息。篇幅所限，这里省略了大部分输出。

不过根据现有的输出内容我们不难发现，`air` 可以使用 `-c` 参数指定配置文件（默认读取执行命令的当前目录下 `.air.toml` 配置文件）；`-v` 参数可以输出版本；`-build.args_bin` 的作用实际上与配置文件中 `[build]` 分类下 `args_bin` 配置项相同。

根据帮助信息的输出内容我们可以总结出，实际上 `air` 所支持的命令就是配置文件中所有的配置项对应的功能。所以学习完 `air` 的配置文件，再来看 `air` 的命令行参数几乎没有学习成本。

那么现在，我们还有最后一个问题需要解决，既然 `air` 即支持命令行参数，又支持配置文件，那么如果同时设置二者，谁的优先级更高呢？

我们一起来做一个实验，修改 Go 程序的 `main.go` 文件，在 `main` 函数的第一行加上如下代码。

```go
fmt.Printf("args[1]: %v\n", os.Args[1])
```

这行代码可以打印 Go 程序接收到的第一个命令行参数。

执行 `air -build.args_bin xxx` 命令启动程序。

![image-20250418173252363](http://images.liangning7.cn/typora/202504181732473.png)

在输出日志中可以看到，我们顺利得到了命令行参数 `xxx`。

接下来我们修改 `.air.toml` 中如下配置项。

```bash
[build]
  args_bin = ["yyy"]
```

现在我们执行 `air -build.args_bin xxx -c .air.toml` 命令启动程序。

![image-20250418173403147](http://images.liangning7.cn/typora/202504181734259.png)

这一次，我们既指定了命令行参数 `-build.args_bin xxx`，又指定了配置文件 `-c .air.toml`。根据输出日志可以发现，`air` 的命令行参数优先级高于配置文件。

这一结论符合直觉，也符合绝大多数 Go 命令行程序的设定，我们可以作为经验记下来。

那么现在如果直接使用 `air -c .air.toml` 命令启动程序。

![image-20250418173453690](http://images.liangning7.cn/typora/202504181734804.png)

输出结果已经不言自明了。

对于 `air` 命令行参数的讲解就到这里，其他参数读者可以自行尝试。

## 总结

在开发阶段使用热重载技术，可以提高我们的开发效率和提升开发体验。

air 是使用 Go 语言编写的一款强大的热重载工具，使用起来非常便捷，只需要一个简单的 `air` 命令即可启动。air 支持通过配置文件或命令行参数来开启或关闭某项功能。推荐使用 air 来热重载你的 Go 应用程序
