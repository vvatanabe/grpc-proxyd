package main

import (
	"context"
	"log"
	"net"

	pb "github.com/vvatanabe/grpc-proxyd/_example/proto/echo"
	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[echo] ")
}

func main() {
	port := ":3001"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	srv := grpc.NewServer()
	pb.RegisterEchoServiceServer(srv,
		&echoService{})
	log.Printf("start server on port%s\n", port)
	if err := srv.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}

type echoService struct{}

func (s *echoService) Echo(ctx context.Context,
	req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Message: req.GetMessage()}, nil
}
