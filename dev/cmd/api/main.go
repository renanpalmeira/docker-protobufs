package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/renanpalmeira/docker-protobufs/api/v1/genproto"
	v1 "github.com/renanpalmeira/docker-protobufs/internal/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	// Automatic load environment variables from .env
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Set global zap logger
	logger, _ := zap.NewProduction()
	_ = zap.ReplaceGlobals(logger)

	// Set grpc address
	grpcAddress := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))

	// Set http address
	httpAddress := fmt.Sprintf(
		"%s:%s",
		os.Getenv("HTTP_HOST"),
		os.Getenv("HTTP_PORT"),
	)

	// Start tcp server
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		zap.L().Fatal("tcp-server", zap.Error(err))
	}

	// Setup context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start http server
	go func() {
		opts := []grpc.DialOption{grpc.WithInsecure()}
		grpcMux := runtime.NewServeMux()

		_ = pb.RegisterHealthHandlerFromEndpoint(ctx, grpcMux, grpcAddress, opts)

		zap.L().Info("http-server", zap.String("description", "up and running http server"))

		if err := http.ListenAndServe(httpAddress, grpcMux); err != nil {
			zap.L().Fatal("http-server", zap.Error(err))
		}
	}()

	// Start grpc server

	grpcServer := grpc.NewServer()

	// Create gRPC handlers
	pb.RegisterHealthServer(grpcServer, v1.NewHealth())

	zap.L().Info("grpc-server", zap.String("description", "up and running grpc server"))

	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Error("grpc-server", zap.Error(err))
	}
}
