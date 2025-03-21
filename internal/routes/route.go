package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/niflheimdevs/backend/internal/wire"
)

func Routes(app *wire.Application) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

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

	mux.Get("/error/{code}", app.ErrorHandler.ReturnError)

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.Put("/user/update-info", app.UserHandler.UpdateUserData)
	mux.Post("/user/profile", app.FileHandler.UploadProfilePhoto)
	mux.Delete("/user/profile", app.FileHandler.DeleteProfilePhoto)
	mux.Put("/user/change-phone/send-otp", app.UserHandler.UpdatePhoneSendOTP)
	mux.Put("/user/change-phone/verify", app.UserHandler.ChangePassword)

	mux.Get("/storage/*", app.FileHandler.GetFile)

	return mux
}
