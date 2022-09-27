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
	"my_grpc/proto/hello_http"
	"net"
)

type HelloService struct{
	hello_http.UnimplementedHelloServer
}

func (t *HelloService) SayHello(ctx context.Context, in *hello_http.HelloHTTPRequest) (*hello_http.HelloHTTPResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}
	if appid != "101010" || appkey != "i am key" {
		return nil, grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}
	resp := new(hello_http.HelloHTTPResponse)
	resp.Message = fmt.Sprintf("Hello %s. Token info: appid=%s,appkey=%s", in.Name, appid, appkey)
	return resp, nil
}
func main() {
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

	listen, _ := net.Listen("tcp", ":9000")
	gServer := grpc.NewServer(grpc.Creds(creds))
	hello_http.RegisterHelloServer(gServer, &HelloService{})
	fmt.Println("Listen ont 9000 with TLS + Token")
	_ = gServer.Serve(listen)
}










