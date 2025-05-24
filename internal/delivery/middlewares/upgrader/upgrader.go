package midupgrader

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/niflheimdevs/backend/bootstrap"
)

type WebSocketUpgrader struct {
	Constants *bootstrap.Constants
}

func NewWebSocketUpgrader(constants *bootstrap.Constants) *WebSocketUpgrader {
	return &WebSocketUpgrader{
		Constants: constants,
	}
}

func (wsUpgrader *WebSocketUpgrader) Upgrade(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
		}

		//lint:ignore SA1029 needed for compatibility with legacy code
		ctx := context.WithValue(r.Context(), wsUpgrader.Constants.Context.WebSocketConnection, conn)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
