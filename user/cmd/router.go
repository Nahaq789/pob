package main

import (
	"pob/user/internal/handler"
	"pob/user/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CreateUserRouter(r gin.IRouter, handler *handler.Userhandler) {
	{
		user := r.Group("/user")
		user.Use(middleware.TraceMiddleware())
		user.POST("/register", handler.Registration)
	}
}

func CreateAuthRouter(r gin.IRouter, handler *handler.AuthHandler) {
	{
		auth := r.Group("/auth")
		auth.Use(middleware.TraceMiddleware())
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)
	}
}
