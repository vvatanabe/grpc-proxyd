package extras

import (
	"context"
	"fmt"
	"strings"

	"github.com/vvatanabe/grpc-proxyd/internal/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func GetDirector(config Config) func(context.Context, string) (context.Context, *grpc.ClientConn, error) {

	credentialsCache := make(map[string]credentials.TransportCredentials)

	return func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		for _, backend := range config.Backends {
			if strings.HasPrefix(fullMethodName, backend.Filter) {
				if config.Verbose {
					fmt.Printf("Found: %s > %s \n", fullMethodName, backend.Backend)
				}
				if backend.CertFile == "" {
					conn, err := grpc.DialContext(ctx, backend.Backend, grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.NewCodec())),
						grpc.WithInsecure())
					return ctx, conn, err
				}
				creds := GetCredentials(credentialsCache, backend)
				if creds != nil {
					conn, err := grpc.DialContext(ctx, backend.Backend, grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.NewCodec())),
						grpc.WithTransportCredentials(creds))
					return ctx, conn, err
				}
				grpclog.Fatalf("Failed to create TLS credentials")
				return ctx, nil, status.Errorf(codes.FailedPrecondition, "Backend TLS is not configured properly in grpc-proxy")
			}
		}
		if config.Verbose {
			fmt.Println("Not found: ", fullMethodName)
		}
		return ctx, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}
}

func GetCredentials(cache map[string]credentials.TransportCredentials, backend Backend) credentials.TransportCredentials {
	if cache[backend.Backend] != nil {
		return cache[backend.Backend]
	}
	creds, err := credentials.NewClientTLSFromFile(backend.CertFile, backend.ServerName)
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
		return nil
	}
	cache[backend.Backend] = creds
	return creds
}
