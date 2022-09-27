package utils

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func GetTlsServerOpts() (opts []grpc.ServerOption) {
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
	return opts
}

func GetTlsDialOpts() (opts []grpc.DialOption) {
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
	return opts
}
