package main

import (
	"context"
	"log/slog"
	"net"
	"pob/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	logger.InitLogger()

	lis, err := net.Listen("tcp", ":9093")
	if err != nil {
		return
	}

	// gRPCサーバー起動
	go func() {
		grpcServer := grpc.NewServer()
		if err := grpcServer.Serve(lis); err != nil {
			slog.ErrorContext(ctx, "gRPC server error", slog.Any("error", err))
		}
	}()

	// webサーバー起動
	router := gin.Default()
	router.Run(":8083")
}

