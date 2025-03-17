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
		var userID int
		authHeader := r.Header.Get("access_token")
		if authHeader == "" {
			userID = -2
		} else {
			claims := am.JWTService.VerifyToken(authHeader)
			if claims == nil {
				userID = -1
			} else {

				userID = int(claims["sub"].(float64))
			}
		}
		ctx := context.WithValue(r.Context(), am.Constants.Context.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
