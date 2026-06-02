package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	"pob/box/cmd/di"
	"pob/box/internal/shared"
	"pob/box/proto"
	"pob/pkg/interceptor/hmac"
	"pob/pkg/logger"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	logger.InitLogger()

	dbClient, err := shared.InitDbClient(ctx, os.Getenv("DB_DSN"))
	if err != nil {
		slog.ErrorContext(ctx, "failed to init db client", slog.Any("error", err))
		return
	}
	defer dbClient.Close()

	publicKeyData, err := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))
	if err != nil {
		slog.ErrorContext(ctx, "failed to read public key", slog.Any("error", err))
		return
	}
	publicKey, err := jwtlib.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse public key", slog.Any("error", err))
		return
	}

	dexConn, err := grpc.NewClient(
		os.Getenv("DEX_GRPC_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(hmac.HmacClientInterceptor(os.Getenv("HMAC_SECRET"))),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to dex", slog.Any("error", err))
		return
	}
	defer dexConn.Close()
	dexClient := proto.NewDexServiceClient(dexConn)

	container, err := di.NewContainer(dbClient, dexClient)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init container", slog.Any("error", err))
		return
	}

	lis, err := net.Listen("tcp", ":9093")
	if err != nil {
		slog.ErrorContext(ctx, "failed to listen", slog.Any("error", err))
		return
	}

	go func() {
		grpcServer := grpc.NewServer()
		proto.RegisterBoxServiceServer(grpcServer, container.Grpc)
		if err := grpcServer.Serve(lis); err != nil {
			slog.ErrorContext(ctx, "gRPC server error", slog.Any("error", err))
		}
	}()

	router := SetupRouter(container, publicKey)
	router.Run(":8083")
}
