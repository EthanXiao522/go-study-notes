# 1.初始化项目：
go mod init go-grpc

# 2.安装工具：

# 3.实现一个gRPC调用案例：
- 编写服务端 .proto 文件
- 生成服务端 .pb.go 文件并同步给客户端
- 编写服务端提供接口的代码
- 编写客户端调用接口的代码
- 启动服务端监听，等待调用
- 启动客户端，调用gRPC

# 4.目录结构：

├─ hello  -- 代码根目录
│  ├─ go_client
│     ├── main.go
│     ├── proto
│         ├── hello
│            ├── hello.pb.go
│  ├─ go_server
│     ├── main.go
│     ├── controller
│         ├── hello_controller
│            ├── hello_server.go
│     ├── proto
│         ├── hello
│            ├── hello.pb.go
│            ├── hello.proto

# 5.代码实现
## 5.1. 编写服务端 hello.proto 文件
```
syntax = "proto3";  //指定 proto 版本：3
package hello;  //包名

//定义服务：service
service Hello {
    //方法1:单项 RPC(一次请求对应一次响应)
    rpc SayHello(HelloReq) returns(HelloRes){}

    //方法2：服务端流式RPC（一次请求，多次响应：流式响应）
    rpc StreamReplies(HelloReq) returns(stream HelloRes){}

    // 客户端流式 RPC
    rpc StreamRequests (stream HelloReq) returns (HelloRes) {}

    //todo：其他RPC待实现
}

//HelloReq 请求结构体
message HelloReq {
    string request = 1; //从1开始，顺序增加
    int32 id = 2;
    string email = 3;

}

//HelloRes 响应结构体
message HelloRes {
    string responce = 1 ;
    int64 timestamp = 2;
}
```

## 5.2. 根据 *.proto生成 *.pb.go文件

```
# 检查 protoc
protoc --version

# 安装插件（确保 $GOBIN 在 PATH）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码（在模块根或以 proto 路径为准）
# 在 go-grpc/go-grpc-server 下生成
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/hello/hello.proto

# 在模块根目录添加依赖并 tidy
cd go-grpc
go get google.golang.org/grpc@latest google.golang.org/protobuf@latest
go mod tidy

# 验证 proto 包能编译
cd go-grpc/go-grpc-server
go build ./proto/...
```

## 5.3. 同时将生成的 hello.pb.go 复制到客户端一份。