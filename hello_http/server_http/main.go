package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	"log"
	gw "my_grpc/proto/hello_http"
	"my_grpc/utils"
	"net/http"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// grpc服务地址
	endpoint := ":9000"
	mux := runtime.NewServeMux()
	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	opts := utils.GetTlsDialOpts()

	// http转grpc
	err := gw.RegisterHelloHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalln("gw.RegisterHelloHandlerFromEndpoint err:", err)
	}
	grpclog.Errorln("HTTP listen on 8000")
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatalln("http.ListenAndServe err:", err)
	}
}






