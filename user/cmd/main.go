package main

import (
	"context"
	"pob/user/internal/handler"
	"pob/user/internal/repository"
	"pob/user/internal/service"
	"pob/user/internal/shared"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	logger := shared.InitLogger()
	dbClient, err := shared.InitDbClient(ctx, logger, "host=localhost user=pob password=pob dbname=user_db port=5432 sslmode=disable")
	if err != nil {
		logger.ErrorContext(ctx, "failed to init db client", "error", err)
		return
	}

	userRepository := repository.NewUserRepository(dbClient)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	CreateUserRouter(router, userHandler)

	router.Run(":8080")
}
