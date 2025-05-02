package midauth

import (
	"context"
	"net/http"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
)

type Authentication struct {
	Constants  *bootstrap.Constants
	JWTService services.JWT
}

func NewAuth(
	Constants *bootstrap.Constants,
	JWTService services.JWT,
) *Authentication {
	return &Authentication{
		Constants:  Constants,
		JWTService: JWTService,
	}
}

func (am *Authentication) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID int
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			userID = -2
		} else {
			tokenString := authHeader[7:]
			claims, err := am.JWTService.VerifyToken(tokenString)
			if err != nil {
				userID = -2
			} else if claims == nil {
				userID = -1
			} else {
				userID = int(claims["sub"].(float64))
			}
		}
		ctx := context.WithValue(r.Context(), am.Constants.Context.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
