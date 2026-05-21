package main

import (
	"context"
	"fmt"
	"go-grpc/go-grpc-server/proto/hello"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 服务端结构体
type Server struct {
	hello.UnimplementedHelloServer
}

// 一元 RPC 实现
func (s *Server) SayHello(ctx context.Context, req *hello.HelloReq) (*hello.HelloRes, error) {
	log.Printf("收到请求：request=%s，id=%d，email=%d", req.Request, req.Id, req.Email)

	if req.Request == "" {
		return nil, status.Error(codes.InvalidArgument, "request 不能为空")
	}

	return &hello.HelloRes{
		Responce:  "收到 Request：" + req.Request,
		Timestamp: time.Now().Unix(),
	}, nil
}

// Stream RPC 实现
func (s *Server) StreamReplies(req *hello.HelloReq, stream hello.Hello_StreamRepliesServer) error {
	log.Printf("收到请求：request=%s，id=%d，email=%d", req.Request, req.Id, req.Email)

	for i := 0; i < 5; i++ {
		resp := &hello.HelloRes{
			Responce:  fmt.Sprintf("响应：%d", i),
			Timestamp: time.Now().Unix(),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

// interceptor
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)

	log.Printf(
		"method=%s ,duration=%s ,err=%v",
		info.FullMethod,
		time.Since(start),
		err,
	)
	return resp, err
}

// 服务端启动监听
func main() {
	lis, err := net.Listen("tcp", ":10001")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	hello.RegisterHelloServer(server, &Server{})

	log.Println("gRPC server start: 10001")

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
