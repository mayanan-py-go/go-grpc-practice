package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"my_grpc/proto/hello_http"
)

const (
	OpenTLS = true
)

// 自定义认证
type customCredential struct {}
// GetRequestMetadata 实现自定义认证接口
func (c *customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
	}, nil
}
// RequireTransportSecurity 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return OpenTLS
}

func main() {
	var opts []grpc.DialOption

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

	// 使用自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	clientConn, err := grpc.Dial(":9000", opts...)
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
