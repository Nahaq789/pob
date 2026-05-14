package main

import (
	"pob/user/internal/handler"

	"github.com/gin-gonic/gin"
)

func CreateUserRouter(r gin.IRouter, handler *handler.Userhandler) {
	{
		user := r.Group("/user")
		user.POST("/register", handler.Registration)
	}
}

func CreateAuthRouter(r gin.IRouter, handler *handler.AuthHandler) {
	{
		auth := r.Group("/auth")
		auth.POST("/login", handler.Login)
	}
}
