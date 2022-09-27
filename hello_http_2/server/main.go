package main

import (
	"context"
	"crypto/tls"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"my_grpc/proto/hello_http"
	"my_grpc/utils"
	"net"
	"net/http"
	"strings"
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
	grpclog.Errorln(in.Name)
	return rsp, nil
}

func main() {
	// 监听
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalln("net.Listen err:", err)
	}

	// grpc TLS Server
	serverOpts := utils.GetTlsServerOpts()
	gServer := grpc.NewServer(serverOpts...)
	hello_http.RegisterHelloServer(gServer, HelloService)

	// gw server
	ctx := context.Background()
	dialOpts := utils.GetTlsDialOpts()
	gwMux := runtime.NewServeMux()
	err = hello_http.RegisterHelloHandlerFromEndpoint(ctx, gwMux, Address, dialOpts)
	if err != nil {
		log.Fatalln("hello_http.RegisterHelloHandlerFromEndpoint err:", err)
	}

	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	srv := http.Server{
		Addr: Address,
		Handler: grpcHandlerFunc(gServer, mux),
		TLSConfig: getTlsConfig(),
	}
	grpclog.Errorf("gRPC and http listen on %s", Address)
	if err = srv.Serve(tls.NewListener(listen, srv.TLSConfig)); err != nil {
		log.Fatalln("srv.Serve err:", err)
	}

}

func getTlsConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("keys/server/server.pem", "keys/server/server.key")
	if err != nil {
		log.Fatalln("tls.LoadX509KeyPair err:", err)
	}
	demoCert := &cert
	return &tls.Config{
		Certificates: []tls.Certificate{*demoCert},
		NextProtos: []string{http2.NextProtoTLS},  // http2 TLS支持
	}
}
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

