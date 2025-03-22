package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	dto "github.com/niflheimdevs/backend/internal/dto/careers"
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

func (gh *GeneralHandler) CreateCareer(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(gh.Constants.Context.UserID).(int)

	params := Validated[dto.PostCareerDTO](gh.Validator, r)

	res := gh.GeneralService.AddCareer(userid, &params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
