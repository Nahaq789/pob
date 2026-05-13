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
