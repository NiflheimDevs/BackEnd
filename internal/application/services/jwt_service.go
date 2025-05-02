package services

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	GenerateToken(userID int) (string, string)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
	RefreshToken(token string) string
}
