package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/niflheimdevs/backend/internal/wire"
)

func Routes(app *wire.Application) http.Handler {
	mux := chi.NewRouter()

	mux.Post("/login", app.UserHandler.Login)

	return mux
}
