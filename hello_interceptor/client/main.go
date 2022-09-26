package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
	"my_grpc/proto/hello"
	"time"
)

const (
	// Address grpc服务地址
	Address = ":9000"

	// OpenTLS 是否开启TLS认证
	OpenTLS = true
)

// 自定义认证
type customCredential struct {}

// GetRequestMetadata 实现自定义认证接口
func (*customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid": "101010",
		"appkey": "i am key",
	}, nil
}

// RequireTransportSecurity 自定义认证是否开启TLS
func (*customCredential) RequireTransportSecurity() bool {
	return OpenTLS
}

func main() {
	var opts []grpc.DialOption

	if OpenTLS {
		// TLS链接
		cert, err := tls.LoadX509KeyPair("keys/client/client.pem", "keys/client/client.key")
		if err != nil {
			log.Fatalln("tls.LoadX509 err:", err)
		}
		certPool := x509.NewCertPool()
		certBytes, err := ioutil.ReadFile("keys/ca/ca.crt")
		if err != nil {
			log.Fatalln("读取ca证书失败：", err)
		}
		certPool.AppendCertsFromPEM(certBytes)
		tcreds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},// 放入客户端证书
			ServerName: "localhost", //证书里面的 commonName
			RootCAs: certPool, // 证书池
		})
		creds := grpc.WithTransportCredentials(tcreds)
		opts = append(opts, creds)
	}else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 指定自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	// 指定客户端interceptor拦截器
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))

	clientConn, err := grpc.Dial(":9000", opts...)
	if err != nil {
		log.Fatalln("grpc.Dial err:", err)
	}
	helloClient := hello.NewHelloClient(clientConn)
	reply, err := helloClient.SayHello(context.Background(), &hello.HelloRequest{Name: "gRPC"})
	if err != nil {
		log.Fatalln("helloClient.SayHello err:", err)
	}
	fmt.Println(reply)
}

// interceptor 客户端拦截器
func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("method=%s req=%v reply=%v duration=%s error=%v\n", method, req, reply, time.Since(start), err)
	return err
}











