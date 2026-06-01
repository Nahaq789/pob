package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	gen "pob/dex/proto"
	"pob/dex/internal/handler"
	"pob/dex/internal/repository"
	"pob/dex/internal/service"
	"pob/dex/internal/shared"
	"pob/pkg/interceptor/hmac"
	"pob/pkg/logger"
	pkgredis "pob/pkg/redis"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	logger.InitLogger()

	dsn := os.Getenv("DB_DSN")
	dbClient, err := shared.InitDbClient(ctx, dsn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init db client", slog.Any("error", err))
		return
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient, err := pkgredis.NewRedisClient(ctx, redisAddr)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init redis client", slog.Any("error", err))
		return
	}
	defer redisClient.Close()

	pokemonRepo := repository.NewPokemonRepository(dbClient, redisClient)
	moveRepo := repository.NewMoveRepository(dbClient)
	abilityRepo := repository.NewAbilityRepository(dbClient)
	itemRepo := repository.NewItemRepository(dbClient, redisClient)

	pokemonSvc := service.NewPokemonService(pokemonRepo)
	moveSvc := service.NewMoveService(moveRepo)
	abilitySvc := service.NewAbilityService(abilityRepo)
	itemSvc := service.NewItemService(itemRepo)

	dexHandler := handler.NewDexHandler(pokemonSvc, moveSvc, abilitySvc, itemSvc)

	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		slog.ErrorContext(ctx, "failed to listen", slog.Any("error", err))
		return
	}

	secret := os.Getenv("HMAC_SECRET")
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(hmac.HmacServerInterceptor(secret)))
	gen.RegisterDexServiceServer(grpcServer, dexHandler)

	slog.InfoContext(ctx, "dex gRPC server listening on :9091")
	if err := grpcServer.Serve(lis); err != nil {
		slog.ErrorContext(ctx, "gRPC server error", slog.Any("error", err))
	}
}
