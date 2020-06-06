package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"

	pb "github.com/vvatanabe/grpc-proxyd/_example/proto/upload"
	"google.golang.org/grpc"
)

const port = ":3003"

func main() {
	log.SetFlags(0)
	log.SetPrefix("[upload] ")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterUploadServiceServer(s, &fileService{})
	log.Printf("start server on port%s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

type fileService struct{}

func (s *fileService) Upload(
	stream pb.UploadService_UploadServer) error {
	var blob []byte
	var name string
	for {
		c, err := stream.Recv()
		if err == io.EOF {
			log.Printf("done %d bytes\n", len(blob))
			break
		}
		if err != nil {
			panic(err)
		}
		name = c.GetName()
		blob = append(blob, c.GetData()...)
	}
	fp := filepath.Join("./upload/resource", name)
	ioutil.WriteFile(fp, blob, 0644)
	stream.SendAndClose(&pb.UploadResponse{
		Size: int64(len(blob))})
	return nil
}
