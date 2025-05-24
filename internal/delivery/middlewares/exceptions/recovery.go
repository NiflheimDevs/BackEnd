package panicwall

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
)

type PanicWall struct {
	Constants *bootstrap.Constants
}

func NewPanicWall(constants *bootstrap.Constants) *PanicWall {
	return &PanicWall{
		Constants: constants,
	}
}

func (recovery *PanicWall) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if r.Header.Get("Upgrade") == "websocket" {
					log.Println(rec)

					conn, ok := r.Context().Value(recovery.Constants.Context.WebSocketConnection).(*websocket.Conn)
					if ok && conn != nil {
						closeMessage := websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "internal server error")
						conn.WriteMessage(websocket.CloseMessage, closeMessage)
						conn.Close()
					} else {
						log.Printf("WebSocket connection not found in context or invalid type")
					}

					return
				}
				log.Println(rec)
				err, ok := rec.(exceptions.Exception)
				if ok {
					json, status := recovery.handleRecoveredError(&err)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(status)
					w.Write(json)
				} else {
					w.WriteHeader(501)
					errorResponse := map[string]interface{}{
						"error": rec,
					}
					jsonResponse, _ := json.Marshal(errorResponse)
					w.Write(jsonResponse)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (recovery *PanicWall) handleRecoveredError(err *exceptions.Exception) ([]byte, int) {
	var code int
	if err.Tag == exceptions.CONFLICT_ERROR {
		code = 409
	} else if err.Tag == exceptions.INTERNAL_ERROR {
		code = 500
	} else if err.Tag == exceptions.NOT_FOUND {
		code = 404
	} else if err.Tag == exceptions.BAD_REQUEST {
		code = 400
	} else if err.Tag == exceptions.LIMIT_EXCEED {
		code = 429
	} else if err.Tag == exceptions.UNAUTHORIZED {
		code = 401
	} else if err.Tag == exceptions.FORBIDDEN {
		code = 403
	} else if err.Tag == exceptions.UNPROCESSABLE {
		code = 422
	} else {
		code = 418
	}

	json, _ := json.Marshal(err)
	return json, code
}
