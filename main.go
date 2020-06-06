package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/vvatanabe/grpc-proxyd/extras"
	"github.com/vvatanabe/grpc-proxyd/internal/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {
	configurationFile := "./config.json"

	args := os.Args[1:]
	if len(args) > 0 {
		configurationFile = args[0]
	}

	config := extras.GetConfiguration(configurationFile)

	listen := ":50051"
	if config.Listen != "" {
		listen = config.Listen
	}

	lis, err := net.Listen("tcp", listen)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("Proxy running at %q\n", listen)

	server := GetServer(config)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func GetServer(config extras.Config) *grpc.Server {
	var opts []grpc.ServerOption

	opts = append(opts, grpc.CustomCodec(proxy.NewCodec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(extras.GetDirector(config))))

	if config.CertFile != "" && config.KeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	return grpc.NewServer(opts...)
}
