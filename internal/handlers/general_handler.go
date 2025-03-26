package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/dto"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/services"
)

type GeneralHandler struct {
	GeneralService *services.GeneralService
	JWTService     *services.JWT
	Constants      *bootstrap.Constants
	Validator      *validator.Validate
}

func NewGeneralHandler(
	generalservice *services.GeneralService,
	jwtService *services.JWT,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *GeneralHandler {
	return &GeneralHandler{
		GeneralService: generalservice,
		JWTService:     jwtService,
		Constants:      constants,
		Validator:      validator,
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

func (gh *GeneralHandler) GetLabel(w http.ResponseWriter, r *http.Request) {
	labels := gh.GeneralService.GetLabels()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(labels); err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.CAST_ERROR},
		})
	}
}

func (gh *GeneralHandler) UpdateCareer(w http.ResponseWriter, r *http.Request) {
	type PutCareerDTO struct {
		Careers []dto.CareerDTO `json:"careers" validate:"required,dive"`
	}
	userid := r.Context().Value(gh.Constants.Context.UserID).(int)

	params := Validated[PutCareerDTO](gh.Validator, r)

	res := gh.GeneralService.UpdateCareer(userid, params.Careers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (gh *GeneralHandler) UpdateUserTag(w http.ResponseWriter, r *http.Request) {
	type PutTags struct {
		Tags []dto.RecieveTagDTO `json:"careers" validate:"required,dive"`
	}

	userid := r.Context().Value(gh.Constants.Context.UserID).(int)

	params := Validated[PutTags](gh.Validator, r)

	res := gh.GeneralService.UpdateTagsForCareerOrUser(userid, params.Tags, true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
