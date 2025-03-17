package panicwall

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type PanicWall struct {
}

func NewPanicWall() *PanicWall {
	return &PanicWall{}
}

func (recovery *PanicWall) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(exceptions.Exception)
				if ok {
					json, status := recovery.handleRecoveredError(&err)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(status)
					w.Write(json)
				} else {
					w.WriteHeader(501)
				}
				log.Println(err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (recovery *PanicWall) handleRecoveredError(err *exceptions.Exception) ([]byte, int) {
	var code int
	if err.Tag == enums.VALIDATION_ERROR {
		code = 409
	} else if err.Tag == enums.AUTHENTICATION_ERROR {
		code = 403
	} else if err.Tag == enums.INTERNAL_ERROR {
		code = 500
	} else if err.Tag == enums.NOT_FOUND {
		code = 404
	} else if err.Tag == enums.BAD_REQUEST {
		code = 400
	} else if err.Tag == enums.LIMIT_EXCEED {
		code = 429
	} else if err.Tag == enums.UNAUTHORIZED {
		code = 401
	} else {
		code = 418
	}

	json, _ := json.Marshal(err)
	return json, code
}
