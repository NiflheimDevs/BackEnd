package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/niflheimdevs/backend/internal/wire"
)

func Routes(app *wire.Application) http.Handler {
	mux := chi.NewRouter()

	mux.Post("/signup/otp", app.UserHandler.ReserveInfo)
	mux.Post("/signup/verify", app.UserHandler.VerifyOTP)
	mux.Get("/redis-test", app.UserHandler.RedisTest)
	return mux
}
