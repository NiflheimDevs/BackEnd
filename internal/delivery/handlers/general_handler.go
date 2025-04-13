package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
)

type GeneralHandler struct {
	TagService    services.TagService
	CareerService services.CareerService
	LabelService  services.LabelService
	JWTService    services.JWT
	Constants     *bootstrap.Constants
	Validator     *validator.Validate
}

func NewGeneralHandler(
	tagService services.TagService,
	careerService services.CareerService,
	labelService services.LabelService,
	jwtService services.JWT,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *GeneralHandler {
	return &GeneralHandler{
		TagService:    tagService,
		CareerService: careerService,
		LabelService:  labelService,
		JWTService:    jwtService,
		Constants:     constants,
		Validator:     validator,
	}
}

func (handler *GeneralHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags := handler.TagService.GetTags()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(tags); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (gh *GeneralHandler) GetLabel(w http.ResponseWriter, r *http.Request) {
	labels := gh.LabelService.GetLabels()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(labels); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (gh *GeneralHandler) UpdateCareer(w http.ResponseWriter, r *http.Request) {
	type PutCareerDTO struct {
		Careers []dto.CareerDTO `json:"careers" validate:"required,dive"`
	}
	userid := r.Context().Value(gh.Constants.Context.UserID).(int)

	params := Validated[PutCareerDTO](gh.Validator, r)

	res := gh.CareerService.UpdateCareers(userid, params.Careers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (gh *GeneralHandler) UpdateUserTag(w http.ResponseWriter, r *http.Request) {
	type PutTags struct {
		Tags []dto.RecieveTagDTO `json:"tags" validate:"required,dive"`
	}

	userid := r.Context().Value(gh.Constants.Context.UserID).(int)

	params := Validated[PutTags](gh.Validator, r)

	res := gh.TagService.UpdateTagsForCareerOrUser(userid, params.Tags, true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
