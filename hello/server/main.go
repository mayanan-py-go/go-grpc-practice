package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"my_grpc/proto/hello"
	"net"
)

const (
	Address = ":9000"
)
type helloService struct {}
var HelloService = new(helloService)
func (*helloService) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	rsp := new(hello.HelloResponse)
	rsp.Message = "hello " + in.Name
	return rsp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalln("net.Listen err:", err)
	}
	gServer := grpc.NewServer()

	// 注册服务
	hello.RegisterHelloServer(gServer, HelloService)

	grpclog.Errorln("Listen on ", Address)
	err = gServer.Serve(listen)
	if err != nil {
		log.Fatalln("gServer.Serve err:", err)
	}
}








