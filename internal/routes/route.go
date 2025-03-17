package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/niflheimdevs/backend/internal/wire"
)

func Routes(app *wire.Application) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			log.Printf("Response Headers: %v", w.Header())
		})
	})
	mux.Use(app.Middlewares.Recovery.Recovery)
	mux.Use(app.Middlewares.RateLimit.RateLimitMiddleware)
	mux.Use(app.Middlewares.Authentication.AuthRequired)

	mux.Post("/signup/send-otp", app.UserHandler.ReserveInfo)
	mux.Post("/signup/verify", app.UserHandler.SignupWithOtp)

	mux.Post("/forget-password/send-otp", app.UserHandler.SendOTP)
	mux.Post("/forget-password/verify", app.UserHandler.VerifyOTP)
	mux.Post("/forget-password/reset", app.UserHandler.ForgetPassword)

	mux.Post("/login", app.UserHandler.Login)
	mux.Post("/change-password", app.UserHandler.ChangePassword)

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return mux
}
