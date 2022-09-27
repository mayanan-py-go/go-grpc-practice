package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"my_grpc/proto/hello_http"
	"my_grpc/utils"
)

func main() {
	// TLS
	opts := utils.GetTlsDialOpts()

	clientConn, err := grpc.Dial(":9000", opts...)
	if err != nil {
		log.Fatalln("grpc.Dial err:", err)
	}
	helloClient := hello_http.NewHelloClient(clientConn)
	reply, err := helloClient.SayHello(context.Background(), &hello_http.HelloHTTPRequest{Name: "gRPC 你好"})
	if err != nil {
		log.Fatalln("helloClient.SayHello err:", err)
	}
	fmt.Println(reply)
}









