package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/niflheimdevs/backend/wire"
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
	//mux.Use(app.Middlewares.RateLimit.RateLimitMiddleware)
	mux.Use(app.Middlewares.Authentication.AuthRequired)

	mux.Post("/login", app.Handlers.UserHandler.Login)

	mux.Post("/signup/send-otp", app.Handlers.UserHandler.ReserveInfo)
	mux.Post("/signup/verify", app.Handlers.UserHandler.SignupWithOtp)

	mux.Post("/forget-password/send-otp", app.Handlers.UserHandler.SendOTP)
	mux.Post("/forget-password/verify", app.Handlers.UserHandler.VerifyOTP)
	mux.Post("/forget-password/reset", app.Handlers.UserHandler.ForgetPassword)

	mux.Post("/change-password", app.Handlers.UserHandler.ChangePassword)

	mux.Get("/tags", app.Handlers.GeneralHandler.GetTags)
	mux.Get("/labels", app.Handlers.GeneralHandler.GetLabel)

	mux.Put("/user/career", app.Handlers.GeneralHandler.UpdateCareer)
	mux.Put("/user/tag", app.Handlers.GeneralHandler.UpdateUserTag)

	mux.Get("/project/user/{user_id}", app.Handlers.ProjectHandler.GetUserProject)
	mux.Get("/project/{project_id}", app.Handlers.ProjectHandler.GetSpeceficProject)
	mux.Post("/project/create", app.Handlers.ProjectHandler.CreateProject)
	mux.Put("/project/{project_id}", app.Handlers.ProjectHandler.UpdateProject)
	mux.Delete("/project/{project_id}", app.Handlers.ProjectHandler.DeleteProject)

	mux.Get("/landing/projects", app.Handlers.ProjectHandler.LandingProps)

	mux.Get("/user/balance", app.Handlers.PaymentHandler.GetUserBalance)
	mux.Get("/transaction", app.Handlers.PaymentHandler.GetUserTransactions)

	mux.Post("/transaction/deposit", app.Handlers.PaymentHandler.Deposit)
	mux.Post("/transaction/withdraw", app.Handlers.PaymentHandler.Withdraw)

	mux.Get("/user/{id}", app.Handlers.UserHandler.GetUserInfo)
	mux.Put("/user/update-info", app.Handlers.UserHandler.UpdateUserData)
	mux.Put("/user/update-username", app.Handlers.UserHandler.UpdateUsername)
	mux.Put("/user/update-email", app.Handlers.UserHandler.UpdateEmail)
	mux.Put("/user/update-phone/send-otp", app.Handlers.UserHandler.UpdatePhoneSendOTP)
	mux.Put("/user/update-phone/verify", app.Handlers.UserHandler.UpdatePhoneVerify)

	mux.Post("/user/profile", app.Handlers.FileHandler.UploadProfilePhoto)
	mux.Delete("/user/profile", app.Handlers.FileHandler.DeleteProfilePhoto)
	mux.Post("/user/resume", app.Handlers.FileHandler.UploadUserResume)
	mux.Delete("/user/resume", app.Handlers.FileHandler.DeleteUserResume)

	mux.Get("/user/profile", app.Handlers.FileHandler.GetProfilePhoto)
	mux.Get("/user/resume", app.Handlers.FileHandler.GetResume)

	mux.Get("/refresh-token", app.Handlers.UserHandler.RefreshToken)

	mux.Post("/bid", app.Handlers.BidHandler.PutBid)
	mux.Post("/bid/{id}/accept", app.Handlers.BidHandler.AcceptBid)
	mux.Put("/bid/{id}", app.Handlers.BidHandler.UpdateBid)

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return mux
}
