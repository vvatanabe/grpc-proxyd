package main

import (
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	pb "github.com/vvatanabe/grpc-proxyd/_example/proto/download"
	"google.golang.org/grpc"
)

const port = ":3002"

func main() {
	log.SetFlags(0)
	log.SetPrefix("[download] ")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	srv := grpc.NewServer()
	pb.RegisterDownloadServiceServer(srv, &downloadService{})
	log.Printf("start server on port%s\n", port)
	if err := srv.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}

type downloadService struct{}

func (s downloadService) Download(req *pb.DownloadRequest,
	stream pb.DownloadService_DownloadServer) error {
	fp := filepath.Join("./download/resource",
		req.GetName())
	fs, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer fs.Close()
	buf := make([]byte, 1000*1024)
	for {
		n, err := fs.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		err = stream.Send(&pb.DownloadResponse{
			Data: buf[:n],
		})
		if err != nil {
			return err
		}
	}
	return nil
}
