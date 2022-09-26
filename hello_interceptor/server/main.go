package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
	"my_grpc/proto/hello"
	"net"
)

const (
	// Address grpc服务地址
	Address = ":9000"
)

// 定义helloService并实现约定的接口
type helloService struct {}

// HelloService Hello服务
var HelloService = new(helloService)

// SayHello 实现Hello服务接口
func (*helloService) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	rsp := new(hello.HelloResponse)
	rsp.Message = "hello: " + in.Name
	return rsp, nil
}

func main() {
	var opts []grpc.ServerOption
	// 使用tls进行加载key pair对
	certifacate, err := tls.LoadX509KeyPair("keys/server/server.pem", "keys/server/server.key")
	if err != nil {
		log.Fatalln("tls加载X509失败,err:",err)
	}
	// 创建证书池
	certPool := x509.NewCertPool()
	// 向证书池中加入证书
	certBytes, err := ioutil.ReadFile("keys/ca/ca.crt")
	if err != nil {
		log.Fatalln("读取ca.crt证书失败，err:", err)
	}
	// 加载证书从pem文件里面
	certPool.AppendCertsFromPEM(certBytes)
	// 创建credentials对象
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certifacate},  // 服务端证书
		ClientAuth: tls.RequireAndVerifyClientCert,  // 需要并且验证客户端证书
		ClientCAs: certPool,  // 客户端证书池
	})
	opts = append(opts, grpc.Creds(creds))

	// 注册interceptor
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	listen, _ := net.Listen("tcp", Address)
	gServer := grpc.NewServer(opts...)
	hello.RegisterHelloServer(gServer, HelloService)
	fmt.Println("Listen on 9000 with TLS + Token + Interceptor")
	_ = gServer.Serve(listen)
}

// interceptor拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}

func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "无token认证信息")
	}
	var (
		appid string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}
	if appid != "101010" || appkey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "Token认识信息无效：appid=%s, appkey=%s", appid, appkey)
	}
	return nil
}









