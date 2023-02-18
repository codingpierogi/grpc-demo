package main

import (
	"context"
	"log"
	"time"

	bm "github.com/codingpierogi/grpc-demo/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	addr = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := bm.NewBookManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	bookId := "101"
	r, err := c.GetBook(ctx, &wrapperspb.StringValue{Value: bookId})
	if err != nil {
		log.Fatalf("Could not find book: %s", bookId)
	}
	log.Printf("Found %s by %s", r.Title, r.Authors[0])
}
