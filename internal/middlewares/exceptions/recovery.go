package middlewareexception

import (
	"encoding/json"
	"net/http"

	"github.com/niflheimdevs/backend/internal/exceptions"
)

type RecoveryMiddleware struct {
}

func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

func (recovery RecoveryMiddleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					json, status := recovery.handleRecoveredError(err)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(status)
					w.Write(json)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (recovery RecoveryMiddleware) handleRecoveredError(err error) ([]byte, int) {
	if loginError, ok := err.(exceptions.LoginError); ok {
		return handleLoginError(loginError)
	} else if RegistrationError, ok := err.(exceptions.RegistrationError); ok {
		return handleRegistrationError(RegistrationError)
	}
}

func handleRegistrationError(registrationErrors exceptions.RegistrationError) ([]byte, int) {
	errorMessages := make(map[string]map[string]string)
	for _, registrationError := range registrationErrors.FieldErrors() {
		if _, ok := errorMessages[registrationError.Field]; !ok {
			errorMessages[registrationError.Field] = make(map[string]string)
		}
		errorMessages[registrationError.Field][registrationError.Tag] = message
	}

	json, _ := json.Marshal(errorMessages)

	return json, 422
}

func handleLoginError(err exceptions.LoginError) ([]byte, int) {
	return []byte(err.Err), 401
}
