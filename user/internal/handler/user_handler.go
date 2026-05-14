package handler

import (
	"net/http"
	"pob/user/internal/service"
	"pob/user/internal/service/dto/user"

	"github.com/gin-gonic/gin"
)

type Userhandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *Userhandler {
	return &Userhandler{
		userService: s,
	}
}

func (u *Userhandler) Registration(ctx *gin.Context) {
	var user user.UserRegistration
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	rCtx := ctx.Request.Context()
	err := u.userService.Registration(rCtx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "user registered",
	})
}
