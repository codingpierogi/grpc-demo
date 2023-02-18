package main

import (
	"context"
	"log"
	"net"

	bm "github.com/codingpierogi/grpc-demo/protos"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

type server struct {
	bookMap map[string]*bm.Book
	bm.UnimplementedBookManagementServer
}

func (s *server) AddBook(ctx context.Context, in *bm.Book) (*wrappers.StringValue, error) {
	_, exists := s.bookMap[in.Id]
	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "Book already exists. : ", in.Id)
	}
	s.bookMap[in.Id] = in
	return &wrappers.StringValue{Value: in.Id}, nil
}

func (s *server) DeleteBook(ctx context.Context, bookId *wrappers.StringValue) (*bm.Book, error) {
	book, exists := s.bookMap[bookId.Value]
	if exists {
		return book, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Book does not exist. : ", bookId)
}

func (s *server) GetBook(ctx context.Context, bookId *wrappers.StringValue) (*bm.Book, error) {
	book, exists := s.bookMap[bookId.Value]
	if exists {
		return book, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Book does not exist. : ", bookId)
}

func main() {
	server := &server{
		bookMap: make(map[string]*bm.Book),
	}
	initSampleData(server)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bm.RegisterBookManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initSampleData(s *server) {
	s.bookMap["101"] = &bm.Book{Id: "101", Title: "Fall, Or Dodge In Hell", Authors: []string{"Neal Stephenson"}}
	s.bookMap["102"] = &bm.Book{Id: "102", Title: "Dune", Authors: []string{"Frank Herbert"}}
	s.bookMap["103"] = &bm.Book{Id: "103", Title: "The Martian", Authors: []string{"Andy Weir"}}
}
