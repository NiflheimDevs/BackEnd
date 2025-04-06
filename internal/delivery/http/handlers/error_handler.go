package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (errorHandler *ErrorHandler) ReturnError(w http.ResponseWriter, r *http.Request) {
	errorCode := chi.URLParam(r, "code")
	if errorCode == "" {
		http.Error(w, "Error code not provided", http.StatusBadRequest)
		return
	}

	code, _ := strconv.Atoi(errorCode)

	http.Error(w, http.StatusText(code), code)
}
