package main

import (
	"context"
	"log/slog"
	"os"
	"pob/pkg/logger"
	"pob/user/internal/handler"
	"pob/user/internal/repository"
	"pob/user/internal/service"
	"pob/user/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	logger.InitLogger()
	dbClient, err := shared.InitDbClient(ctx, "host=localhost user=pob password=pob dbname=user_db port=5432 sslmode=disable")
	if err != nil {
		slog.ErrorContext(ctx, "failed to init db client", "error", err)
		return
	}

	// user
	userRepository := repository.NewUserRepository(dbClient)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	// auth
	keyPath := os.Getenv("PRIVATE_KEY_PATH")
	keyData, ioErr := os.ReadFile(keyPath)
	if ioErr != nil {
		slog.ErrorContext(ctx, "failed to read private key path", "error", ioErr)
		return
	}
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	authRepository := repository.NewAuthRepository(dbClient, privateKey)
	refreshTokenRepository := repository.NewRefreshTokenRepository(dbClient)
	authService := service.NewAuthService(authRepository, refreshTokenRepository)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	CreateUserRouter(router, userHandler)
	CreateAuthRouter(router, authHandler)

	router.Run(":8080")
}
