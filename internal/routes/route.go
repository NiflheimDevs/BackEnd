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
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request Headers: %v", r.Header)
			next.ServeHTTP(w, r)
			log.Printf("Response Headers: %v", w.Header())
		})
	})

	mux.Use(app.Middlewares.Recovery.Recovery)
	mux.Use(app.Middlewares.RateLimit.RateLimitMiddleware)
	mux.Use(app.Middlewares.Authentication.AuthRequired)

	mux.Post("/login", app.UserHandler.Login)

	mux.Post("/signup/send-otp", app.UserHandler.ReserveInfo)
	mux.Post("/signup/verify", app.UserHandler.SignupWithOtp)

	mux.Post("/forget-password/send-otp", app.UserHandler.SendOTP)
	mux.Post("/forget-password/verify", app.UserHandler.VerifyOTP)
	mux.Post("/forget-password/reset", app.UserHandler.ForgetPassword)

	mux.Post("/change-password", app.UserHandler.ChangePassword)

	mux.Get("/tags", app.GeneralHandler.GetTags)

	mux.Put("/user/career", app.GeneralHandler.UpdateCareer)
	mux.Put("/user/tag", app.GeneralHandler.UpdateUserTag)

	mux.Get("/project", app.ProjectHandler.GetUserProject)
	mux.Get("/project/{project_id}", app.ProjectHandler.GetSpeceficProject)

	mux.Post("/project/create", app.ProjectHandler.CreateProject)
	mux.Put("/project/{project_id}", app.ProjectHandler.UpdateProject)
	mux.Delete("/project/{project_id}", app.ProjectHandler.DeleteProject)

	mux.Get("/error/{code}", app.ErrorHandler.ReturnError)

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.Get("/user/{id}", app.UserHandler.GetUserInfo)
	mux.Put("/user/update-info", app.UserHandler.UpdateUserData)
	mux.Put("/user/update-username", app.UserHandler.UpdateUsername)
	mux.Put("/user/update-email", app.UserHandler.UpdateEmail)
	mux.Put("/user/update-phone/send-otp", app.UserHandler.UpdatePhoneSendOTP)
	mux.Put("/user/update-phone/verify", app.UserHandler.UpdatePhoneVerify)

	mux.Post("/user/profile", app.FileHandler.UploadProfilePhoto)
	mux.Delete("/user/profile", app.FileHandler.DeleteProfilePhoto)

	mux.Get("/storage/*", app.FileHandler.GetFile)

	return mux
}
