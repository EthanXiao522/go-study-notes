package main

import (
	"context"
	"go-grpc/go-grpc-client/proto/hello"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	/* 	conn, err := grpc.Dial("localhost:10001", //连接服务端ip:port
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 不安全的连接（仅开发用）
	) */
	conn, err := grpc.NewClient("localhost:10001", //连接服务端ip:port
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 不安全的连接（仅开发用）
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	client := hello.NewHelloClient(conn)

	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// metadata
	md := metadata.New(map[string]string{
		"token": "abc123",
	})

	ctx = metadata.NewOutgoingContext(ctx, md)

	//一元 RPC
	resp, err := client.SayHello(ctx, &hello.HelloReq{
		Request: "一元 RPC 请求",
		Id:      1,
		Email:   "aaa@123.com",
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("SayHello:", resp)

	log.Println("##############分割线####################")

	//01.启动Gin+gRPC模式(Gin对外提供api，gRPC内部服务调用)
	//127.0.0.1:8080/sayhello?Request=ni hao&id=223&Email=email
	r := gin.Default()
	//02.在路由处理函数中使用 gRPC 客户端
	r.GET("/sayhello", func(ctx *gin.Context) {
		id, _ := strconv.ParseInt(ctx.Query("Request"), 10, 32)
		resp, err := client.SayHello(ctx, &hello.HelloReq{
			Request: ctx.Query("Request"),
			Id:      int32(id),
			Email:   ctx.Query("Email"),
		})
		if err != nil {
			// 错误处理：将 gRPC 错误码转换为 HTTP 状态码[citation:3]
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, resp)
	})

	log.Println("##############分割线####################")

	//stream RPC
	stream, err := client.StreamReplies(
		context.Background(),
		&hello.HelloReq{
			Request: "stream RPC 请求",
			Id:      2,
			Email:   "bbb@123.com",
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	//接收stream响应
	for {
		resp, err := stream.Recv()

		if err == io.EOF { //响应结束
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println("stream response:", resp)
	}

	log.Println("##############分割线####################")

	// 客户端流式 RPC 调用
	clientStreamCall(client)

	r.Run(":8080")
}

// 客户端流式 RPC 调用
func clientStreamCall(client hello.HelloClient) {
	log.Println("\n=== 客户端流式 RPC 调用 ===")

	ctx := context.Background()
	stream, err := client.StreamRequests(ctx)
	if err != nil {
		log.Printf("创建流失败: %v", err)
		return
	}

	names := []string{"王五", "赵六", "小明", "小红"}
	for _, name := range names {
		req := &hello.HelloReq{
			Request: name,
			Id:      100,
			Email:   "bbb@123.com",
		}
		if err := stream.Send(req); err != nil {
			log.Printf("发送失败: %v", err)
			return
		}
		log.Printf("发送: name=%s", name)
		time.Sleep(200 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("接收响应失败: %v", err)
		return
	}

	log.Printf("最终响应: %s", resp.Responce)
}
