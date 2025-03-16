package midauth

import (
	"context"
	"net/http"

	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/services"
)

type Authentication struct {
	Constants  *bootstrap.Constants
	JWTService *services.JWT
}

func NewAuth(
	Constants *bootstrap.Constants,
	JWTService *services.JWT,
) *Authentication {
	return &Authentication{
		Constants:  Constants,
		JWTService: JWTService,
	}
}

func (am *Authentication) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("access_token")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := authHeader
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}
		claims := am.JWTService.VerifyToken(tokenString)

		userID := int(claims["sub"].(float64))
		ctx := context.WithValue(r.Context(), am.Constants.Context.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
