package handler

import (
	"errors"
	"net/http"
	"pob/user/internal/model/apperror"
	"pob/user/internal/service"
	"pob/user/internal/service/dto/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: s,
	}
}

func (a *AuthHandler) Login(ctx *gin.Context) {
	var loginReq auth.LoginRequest
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

	ctx.SetCookie("refresh_token", token.RefreshToken, 60*60*24*7, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  token.AccessToken,
	})
}

func (a *AuthHandler) Refresh(ctx *gin.Context) {
	rt, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  err.Error(),
		})
		return
	}

	rCtx := ctx.Request.Context()
	token, err := a.authService.Refresh(rCtx, rt)
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

	ctx.SetCookie("refresh_token", token.RefreshToken, 60*60*24*7, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  token.AccessToken,
	})
}
