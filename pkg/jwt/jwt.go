package jwt

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenString string, publicKey *rsa.PublicKey) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func ExtractBearerToken(authHeader string) string {
	return strings.TrimPrefix(authHeader, "Bearer ")
}
