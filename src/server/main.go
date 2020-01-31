package main

import (
	"../config"
	pb "./sUrl"
	"crypto/tls"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"./db"
)

// 定义ShortUrlService并实现约定的接口
type shortUrlService struct{}

// ShortUrlService Hello HTTP服务
var ShortUrlService = shortUrlService{}

// GetShortProduct 实现GetShortProduct服务接口
func (h shortUrlService) GetShortProduct(ctx context.Context, in *pb.GetShortProductRequest) (*pb.GetShortProductReponse, error) {
	resp := new(pb.GetShortProductReponse)
	resp.Message = "Hello " + in.SourceUrl + "."
	return resp, nil
}
func init(){
	//初始化mysql连接
	db.Init()
}
func main() {
	config, err := config.NewConfigure()
	endpoint := strings.Join([]string{config.Host, config.Port}, ":")
	conn, err := net.Listen("tcp", endpoint)
	if err != nil {
		grpclog.Fatalf("TCP Listen err:%v\n", err)
	}
	// grpc tls server
	creds, err := credentials.NewServerTLSFromFile("./cert/server.pem", "./cert/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to create server TLS credentials %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterShortUrlServiceServer(grpcServer, ShortUrlService)
	// gw server
	ctx := context.Background()
	dcreds, err := credentials.NewClientTLSFromFile("./cert/server.pem", "test.com")
	if err != nil {
		grpclog.Fatalf("Failed to create client TLS credentials %v", err)
	}
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := runtime.NewServeMux()
	if err = pb.RegisterShortUrlServiceHandlerFromEndpoint(ctx, gwmux, endpoint, dopts); err != nil {
		grpclog.Fatalf("Failed to register gw server: %v\n", err)
	}
	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	srv := &http.Server{
		Addr:      endpoint,
		Handler:   grpcHandlerFunc(grpcServer, mux),
		TLSConfig: getTLSConfig(),
	}
	//项目初始化相关
	this.init()
	grpclog.Infof("gRPC and https listen on: %s\n", endpoint)
	fmt.Printf("gRPC and https listen on: %s\n", endpoint)
	if err = srv.Serve(tls.NewListener(conn, srv.TLSConfig)); err != nil {
		grpclog.Fatal("ListenAndServe: ", err)
	}
	return
}
func getTLSConfig() *tls.Config {
	cert, _ := ioutil.ReadFile("./cert/server.pem")
	key, _ := ioutil.ReadFile("./cert/server.key")
	var demoKeyPair *tls.Certificate
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		grpclog.Fatalf("TLS KeyPair err: %v\n", err)
	}
	demoKeyPair = &pair
	return &tls.Config{
		Certificates: []tls.Certificate{*demoKeyPair},
		NextProtos:   []string{http2.NextProtoTLS}, // HTTP2 TLS支持
	}
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
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
