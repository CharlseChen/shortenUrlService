package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "../server/sUrl"
)

func main() {
	cert, err := tls.LoadX509KeyPair("../server/cert/server.pem", "../server/cert/server.key")
	if err != nil {
		panic(err)
	}

	// 将根证书加入证书池
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("../server/cert/server.pem")
	if err != nil {
		panic(err)
	}

	if !certPool.AppendCertsFromPEM(bs) {
		panic("cc")
	}

	// 新建凭证
	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "test.com",
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})

	dialOpt := grpc.WithTransportCredentials(transportCreds)

	conn, err := grpc.Dial("localhost:50053", dialOpt)
	if err != nil {
		log.Fatalf("Dial failed:%v", err)
	}
	defer conn.Close()

	client := pb.NewShortUrlServiceClient(conn)
	resp1, err := client.GetShortProduct(context.Background(), &pb.GetShortProductRequest{
		SourceUrl: "Hello Server 1 !!",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Resp1:%+v", resp1)


}