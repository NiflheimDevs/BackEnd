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
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
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
	mux.Get("/project/participated/{user_id}", app.Handlers.ProjectHandler.GetParticipatedProjectsForUser)
	mux.Get("/project/{project_id}", app.Handlers.ProjectHandler.GetSpeceficProject)
	mux.Post("/project/create", app.Handlers.ProjectHandler.CreateProject)
	mux.Put("/project/{project_id}", app.Handlers.ProjectHandler.UpdateProject)
	mux.Delete("/project/{project_id}", app.Handlers.ProjectHandler.DeleteProject)
	mux.Post("/project/{project_id}/done", app.Handlers.ProjectHandler.EndProject)

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
	mux.Post("/team/profile/{team_id}", app.Handlers.FileHandler.UploadTeamProfilePhoto)
	mux.Delete("/team/profile/{team_id}", app.Handlers.FileHandler.DeleteTeamProfilePhoto)
	// mux.Get("/user/{id}/project", app.Handlers.ProjectHandler.GetOneManTeamProjects)

	mux.Get("/user/profile", app.Handlers.FileHandler.GetProfilePhoto)
	mux.Get("/user/resume", app.Handlers.FileHandler.GetResume)

	mux.Get("/refresh-token", app.Handlers.UserHandler.RefreshToken)

	mux.Post("/bid", app.Handlers.BidHandler.PutBid)
	mux.Post("/bid/{id}/accept", app.Handlers.BidHandler.AcceptBid)
	mux.Put("/bid/{id}", app.Handlers.BidHandler.UpdateBid)
	mux.Get("/project/{project_id}/bid", app.Handlers.BidHandler.GetProjectBids)
	mux.Get("/project/{project_id}/view", app.Handlers.BidHandler.ViewBidsOfTheProject)

	mux.Post("/team", app.Handlers.TeamHandler.CreateTeam)
	mux.Patch("/team", app.Handlers.TeamHandler.UpdateTeamInfo)
	mux.Delete("/team/{team_id}", app.Handlers.TeamHandler.DeleteTeam)
	mux.Get("/team/user/{user_id}", app.Handlers.TeamHandler.GetTeamsForUser)
	mux.Get("/team/{team_id}", app.Handlers.TeamHandler.GetTeam)
	mux.Patch("/team/member/pos", app.Handlers.TeamHandler.UpdateMemberPosition)
	mux.Patch("/team/member/role", app.Handlers.TeamHandler.UpdateMemberRole)
	mux.Delete("/team/member", app.Handlers.TeamHandler.RemoveMember)
	mux.Post("/team/member", app.Handlers.TeamHandler.AddMembers)
	mux.Get("/team/{team_id}/project", app.Handlers.ProjectHandler.GetTeamProjects)
	mux.Get("/team/bidding", app.Handlers.TeamHandler.GetTeamsForBidding)
	mux.Get("/team/{team_id}/bid", app.Handlers.BidHandler.GetTeamBids)

	mux.Get("/role/team", app.Handlers.RoleHandler.GetTeamRoles)
	mux.Get("/role/team/{role}", app.Handlers.RoleHandler.GetPermissionsForRole)

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return mux
}
