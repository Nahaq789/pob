package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"pob/user/internal/model/apperror"
	"pob/user/internal/service"
	"pob/user/internal/service/dto/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(s *service.AuthService, l *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: s,
		logger:      l,
	}
}

func (a *AuthHandler) Login(ctx *gin.Context) {
	var loginReq auth.Login
	if err := ctx.ShouldBind(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	rCtx := ctx.Request.Context()
	token, err := a.authService.Login(rCtx, loginReq)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  token,
	})
}
