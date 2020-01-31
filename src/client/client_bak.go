package main
import (
	"fmt"
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "../server/sUrl"
)

const (
	address = "localhost:50053"
)
func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewShortUrlServiceClient(conn)
	reply, err := client.GetShortProduct(context.Background(), &pb.GetShortProductRequest{SourceUrl: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetTargetUrl())
	fmt.Println(reply.GetMessage())
}
