package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"my_grpc/proto/hello_http"
	"net"
)

const (
	Address = ":9000"
)
type helloService struct {
	hello_http.UnimplementedHelloServer
}
var HelloService = new(helloService)
func (*helloService) SayHello(ctx context.Context, in *hello_http.HelloHTTPRequest) (*hello_http.HelloHTTPResponse, error) {
	rsp := new(hello_http.HelloHTTPResponse)
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
	hello_http.RegisterHelloServer(gServer, HelloService)

	grpclog.Errorln("Listen on ", Address)
	err = gServer.Serve(listen)
	if err != nil {
		log.Fatalln("gServer.Serve err:", err)
	}
}








