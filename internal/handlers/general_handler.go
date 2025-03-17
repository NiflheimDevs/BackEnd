package handlers

import (
	"net/http"

	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/services"
)

type GeneralHandler struct {
	GeneralService *services.GeneralService
}

func NewGeneralHandler(generalservice *services.GeneralService) *GeneralHandler {
	return &GeneralHandler{
		GeneralService: generalservice,
	}
}

func (handler *GeneralHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags := handler.GeneralService.GetTags()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(tags); err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.CAST_ERROR},
		})
	}
}
