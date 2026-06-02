package middleware

import (
	"crypto/rsa"
	"net/http"
	"pob/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const userIdKey = "user_id"

func AuthMiddleware(publicKey *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		at := jwt.ExtractBearerToken(c.Request.Header.Get("Authorization"))
		claims, err := jwt.VerifyToken(at, publicKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		c.Set(userIdKey, claims.UserID)
		c.Next()
	}
}
