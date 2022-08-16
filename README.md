# gRPC快速入门

 	RPC算是近些年比较火热的概念了，随着微服务架构的兴起，RPC的应用越来越广泛。本文介绍了RPC和gRPC的相关概念，并且通过详细的代码示例介绍了gRPC的基本使用。



## protobuf简介

`Protocol Buffers(protobuf)`：与编程语言无关，与程序运行平台无关的**数据序列化协议**以及**接口定义语言**(IDL: interface definition language)。

要使用`protobuf`需要先理解几个概念：

* `protobuf`编译器`protoc`，用于编译`.proto`文件

  * 开源地址：https://github.com/protocolbuffers/protobuf

* 编程语言的`protobuf`插件，搭配`protoc`编译器，根据`.proto`文件生成对应编程语言的代码。

* `protobuf runtime library`：每个编程语言有各自的`protobuf runtime`，用于实现各自语言的`protobuf`协议。

* Go语言的`protobuf`插件和`runtime library`有过2个版本：

  * 第1个版本开源地址：[https://github.com/golang/protobuf](https://github.com/golang/protobuf)，包含有插件`proto-gen-go`，可以生成`xx.pb.go`和`xx_grpc.pb.go`。Go工程里导入该版本的`protobuf runtime`的方式如下：

    ```go
    import "github.com/golang/protobuf"
    ```

  * 第2个版本开源地址：[https://github.com/protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go)，同样包含有插件`proto-gen-go`。不过该项目的`proto-gen-go`从`v1.20`版本开始，不再支持生成gRPC服务定义，也就是`xx_grpc.pb.go`文件。要生成gRPC服务定义需要使用`grpc-go`里的`progo-gen-go-grpc`插件。Go工程里导入该版本的`protobuf runtime`的方式如下：

    ```go
    import "google.golang.org/protobuf"
    ```

  推荐使用第2个版本，对protobuf的API做了优化和精简，并且把工程界限分清楚了：

  * 第一，把`protobuf`的Go实现都放在protobuf的项目里，而不是放在golang语言项目下面。
  * 第二，把`gRPC`的生成，放在`grpc-go`项目里，而不是和`protobuf runtime`混在一起。

  有的老项目可能使用了第1个版本的`protobuf runtime`，在老项目里开发新功能的时候也可以使用第2个版本`protobuf runtime`，支持2个版本在一个Go项目里共存。但是要**注意**：一个项目里同时使用2个版本必须保证第一个版本的版本号不低于`v1.4`。

## RPC是什么

​	在分布式计算，远程过程调用（英语：Remote Procedure Call，缩写为 RPC）是一个计算机通信协议。该协议允许运行于一台计算机的程序调用另一个地址空间（通常为一个开放网络的一台计算机）的子程序，而程序员就像调用本地程序一样，无需额外地为这个交互作用编程（无需关注细节）。RPC是一种服务器-客户端（Client/Server）模式，经典实现是一个通过`发送请求-接受回应`进行信息交互的系统。



## gRPC是什么

​	`gRPC`是一种现代化开源的高性能RPC框架，能够运行于任意环境之中。最初由谷歌进行开发。它使用HTTP/2作为传输协议。

​	在gRPC里，客户端可以像调用本地方法一样直接调用其他机器上的服务端应用程序的方法，帮助你更容易创建分布式应用程序和服务。与许多RPC系统一样，gRPC是基于定义一个服务，指定一个可以远程调用的带有参数和返回类型的的方法。在服务端程序中实现这个接口并且运行gRPC服务处理客户端调用。在客户端，有一个stub提供和服务端相同的方法。 

![img](https://img2020.cnblogs.com/blog/1328551/202102/1328551-20210226145122616-1863140796.png)

 

## 为什么要用gRPC

​	使用gRPC， 我们可以一次性的在一个`.proto`文件中定义服务并使用任何支持它的语言去实现客户端和服务端，反过来，它们可以应用在各种场景中，从Google的服务器到你自己的平板电脑—— gRPC帮你解决了不同语言及环境间通信的复杂性。使用`protocol buffers`还能获得其他好处，包括高效的序列号，简单的IDL以及容易进行接口更新。总之一句话，使用gRPC能让我们更容易编写跨语言的分布式代码。







## gRPC-Go简介

gRPC-Go: gRPC的Go语言实现，基于HTTP/2的RPC框架。

开源地址：https://github.com/grpc/grpc-go

Go项目里导入该模块的方式如下：

```go
import "google.golang.org/grpc"
```

`grpc-go`项目里还包含有`protoc-gen-go-grpc`插件，用于根据`.proto`文件生成`xx_grpc.pb.go`文件。





## 环境安装

分为3步：

* 安装Go

  * 步骤参考：https://go.dev/doc/install

* 安装Protobuf编译器`protoc`: 用于编译`.proto` 文件

  * 步骤参考：https://grpc.io/docs/protoc-installation/

  * 执行如下命令查看`protoc`的版本号，确认版本号是3+，用于支持protoc3

    ```bash
    protoc --version
    ```

* 安装`protoc`编译器的Go语言插件

  * `protoc-gen-go`插件：用于生成`xx.pb.go`文件

    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    ```

  * `protoc-gen-go-grpc`插件：用于生成`xx_grpc.pb.go`文件

    ```bash
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

**注意**：有的教程可能只让你安装`protoc-gen-go`，没有安装`protoc-gen-go-grpc`，那有2种情况：

* 使用的是第1个版本`github.com/golang/protobuf`的`protoc-gen-go`插件。
* 使用的是第2个版本`google.golang.org/protobuf`的`protoc-gen-go`插件并且`protoc-gen-go`版本号低于`v1.20`。从`v1.20`开始，第2个版本的`protoc-gen-go`插件不再支持生成gRPC服务定义。下面是官方说明：

> The v1.20 [`protoc-gen-go`](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) does not support generating gRPC service definitions. In the future, gRPC service generation will be supported by a new `protoc-gen-go-grpc` plugin provided by the Go gRPC project.
>
> The `github.com/golang/protobuf` version of `protoc-gen-go` continues to support gRPC and will continue to do so for the foreseeable future.





## gRPC入门示例

### 编写proto代码

gRPC是基于Protocol Buffers。

`Protocol Buffers`是一种与语言无关，平台无关的可扩展机制，用于序列化结构化数据。使用`Protocol Buffers`可以一次定义结构化的数据，然后可以使用特殊生成的源代码轻松地在各种数据流中使用各种语言编写和读取结构化数据。

关于`Protocol Buffers`的教程可以自行在网上搜索，本文默认读者熟悉`Protocol Buffers`。

```protobuf
syntax = "proto3"; // 版本声明，使用Protocol Buffers v3版本

option go_package = "hello_server/pb";

package pb; // 包名


// 定义一个打招呼服务
service Greeter {
    // SayHello 方法
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// 包含人名的一个请求消息
message HelloRequest {
    string name = 1;
}

// 包含问候语的响应消息
message HelloReply {
    string message = 1;
}
```

执行下面的命令，生成go语言源代码：

```bash
protoc -I hello_server/ hello_server/pb/hello.proto --go_out=plugins=grpc:hello_server/pb
```

在`gRPC_demo/hello_server/pb`目录下会生成`hello.pb.go`文件。

### 编写Server端Go代码

```go
package main

import (
	"context"
	"fmt"
	"hello_server/pb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务

	reflection.Register(s) //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}

```

将上面的代码保存到`gRPC_demo/hello_server/main.go`文件中，编译并执行：

```bash
go run hello_server/main.go
```

### 编写Client端Go代码

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"hello_client/pb"

	"google.golang.org/grpc"
)

// hello_client

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

```

将上面的代码保存到`gRPC_demo/hello_client/main.go`文件中，编译并执行：

```bash
go run hello_client/main.go
```

得到输出如下（注意要先启动server端再启动client端）：

```bash
$ ./client 
Greeting: Hello q1mi!
```





### gRPC跨语言调用

接下来，我们演示一下如何使用gRPC实现跨语言的RPC调用。

我们使用`Python`语言编写`Client`，然后向上面使用`go`语言编写的`server`发送RPC请求。

### 生成Python代码

在`gRPC_demo`目录执行下面的命令：

```bash
python -m grpc_tools.protoc -I py_client/pb/ --python_out=py_client/ --grpc_python_out=py_client/ client_py/pb/hello.proto
```

上面的命令会在`gRPC_demo/helloworld/client/`目录生成如下两个python文件：

```bash
helloworld_pb2.py
helloworld_pb2_grpc.py
```

### 编写Python版Client

在`gRPC_demo/py_client/`目录中编写`client.py`文件，其内容如下：

```python
from __future__ import print_function

import logging

import grpc
import hello_pb2
import hello_pb2_grpc


def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    with grpc.insecure_channel('127.0.0.1:8972') as channel:
        stub = hello_pb2_grpc.GreeterStub(channel)
        response = stub.SayHello(hello_pb2.HelloRequest(name='q1mi'))
    print("Greeter client received: " + response.message)


if __name__ == '__main__':
    logging.basicConfig()
    run()
```

将上面的代码保存执行，得到输出结果如下：

```bash
gRPC_demo $ python py_client/client.py 
Greeter client received: Hello q1mi!
```

这里我们就实现了，使用python代码编写的client去调用Go语言版本的server了。





## 官方示例

### 下载代码

以`grpc-go`的v1.41.0版本为例，下载代码并进入到`grpc-go/examples/helloworld`目录：

```bash
git clone -b v1.41.0 https://github.com/grpc/grpc-go
cd grpc-go/examples/helloworld
```

### 运行代码

* 启动服务端

  ```bash
  go run greeter_server/main.go
  ```

  终端会打印如下内容，表示服务端已经启动并且在监听`50051`端口

  ```bash
  2022/01/02 13:01:08 server listening at [::]:50051
  ```

* 启动客户端。客户端会发送`SayHello`请求给服务端

  ```bash
  go run greeter_client/main.go
  ```

  终端会打印如下内容，表示收到了服务端的响应。

  ```bash
  2022/01/02 13:01:25 Greeting: Hello world
  ```

  



## 工程开发

自己在使用`protobuf`和`grpc-go`开发的时候，按照如下步骤来操作：

* 定义`.proto`文件，包括消息体和rpc服务接口定义
* 使用`protoc`命令来编译`.proto`文件，用于生成`xx.pb.go`和`xx_grpc.pb.go`文件
* 在服务端实现rpc里定义的方法
* 客户端调用rpc方法，获取响应结果

我们通过对上面的`grpc-go/examples/helloworld`做修改，来说明上述步骤。

* 第一步，在`helloworld.proto`里增加一个rpc方法`SayHelloAgain`，参数和返回值和`SayHello`保持一样。

  ```protobuf
  // The greeting service definition.
  service Greeter {
    // Sends a greeting
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    // send another greeting
    rpc SayHelloAgain (HelloRequest) returns (HelloReply) {}
  }
  ```

* 第二步，在`grpc-go/examples/helloworld`目录使用`protoc`命令编译`.proto`文件，生成新的`helloworld.pb.go`和`helloword_grpc.pb.go`文件。命令如下：

  ```bash
  protoc --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      helloworld/helloworld.proto
  ```

* 第三步，在服务端实现rpc里新定义的方法`SayHelloAgain`。在`greeter_server/main.go`添加如下代码：

  ```go
  func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
  	log.Printf("Received: %v", in.GetName())
  	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
  }
  ```

* 第四步，在客户端调用新定义的rpc方法，获取响应结果。在`greeter_client/main.go`添加如下代码：

  ```go
  r2, err2 := c.SayHelloAgain(ctx, &pb.HelloRequest{Name: *name})
  if err2 != nil {
  	log.Fatalf("could not greet: %v", err2)
  }
  log.Printf("Greeting: %s", r2.GetMessage())
  ```

* 第五步，运行程序

  * 先启动服务端

    ```bash
    go run greeter_server/main.go
    ```

  * 再启动客户端

    ```bash
    go run greeter_client/main.go Alice
    ```

客户端会打印如下内容：

```bash
2022/01/02 13:37:58 Greeting: Hello alice
2022/01/02 13:37:58 Greeting: Hello again alice
```

至此，我们就对如何在Go工程里使用`protobuf`和`gRPC`有了一个初步的了解和入门。



## 进阶学习

想要进一步学习，主要是深入了解`protobuf`和`gRPC`在Go语言里的使用技巧和原理

* `protobuf`官方学习地址：

  * https://developers.google.com/protocol-buffers/docs/proto3
  * https://developers.google.com/protocol-buffers/docs/gotutorial
  * https://developers.google.com/protocol-buffers/docs/reference/go-generated
  * https://developers.google.com/protocol-buffers/docs/reference/proto3-spec

* `gRPC`官方学习地址：

  * https://grpc.io/docs/languages/go/

  



## References

* https://grpc.io/docs/languages/go/quickstart/

* https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.20.0#v1.20-grpc-support
* https://stackoverflow.com/questions/64828054/differences-between-protoc-gen-go-and-protoc-gen-go-grpc
* https://github.com/golang/protobuf
* https://github.com/protocolbuffers/protobuf-go