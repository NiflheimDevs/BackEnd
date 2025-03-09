package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/niflheimdevs/backend/internal/wire"
)

func Routes(app *wire.Application) http.Handler {
	mux := chi.NewRouter()

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

	return mux
}
