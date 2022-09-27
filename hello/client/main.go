package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"my_grpc/proto/hello_http"
)

func main() {
	clientConn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("grpc.Dial err:", err)
	}
	helloClient := hello_http.NewHelloClient(clientConn)
	reply, err := helloClient.SayHello(context.Background(), &hello_http.HelloHTTPRequest{Name: "gRPC"})
	if err != nil {
		log.Fatalln("helloClient.SayHello err:", err)
	}
	fmt.Println(reply)
}









