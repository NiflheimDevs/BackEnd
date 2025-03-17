package services

import (
	"encoding/json"

	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/repositories"
)

type GeneralService struct {
	GeneralRepo *repositories.GeneralRepo
}

func NewGeneralService(GeneralRepo *repositories.GeneralRepo) *GeneralService {
	return &GeneralService{
		GeneralRepo: GeneralRepo,
	}
}

func (service *GeneralService) GetTags() []byte {
	tags, err := service.GeneralRepo.GetTags()
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
		})
	}

	marshaled, err := json.Marshal(tags)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.CAST_ERROR},
		})
	}

	return marshaled
}
