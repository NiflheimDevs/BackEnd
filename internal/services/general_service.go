package services

import (
	"encoding/json"

	dto "github.com/niflheimdevs/backend/internal/dto/careers"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/models"
	"github.com/niflheimdevs/backend/internal/repositories"
)

type GeneralService struct {
	GeneralRepo *repositories.GeneralRepo
	UserRepo    *repositories.UserRepo
}

func NewGeneralService(GeneralRepo *repositories.GeneralRepo, userRepo *repositories.UserRepo) *GeneralService {
	return &GeneralService{
		GeneralRepo: GeneralRepo,
		UserRepo:    userRepo,
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

func (gs *GeneralService) GetCareerForUser(userid int) *dto.SendCareersDTO {
	var res dto.SendCareersDTO
	careers, err := gs.GeneralRepo.GetCareerForUser(userid)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
		})
	}

	res.Careers = careers
	var tags []models.TagModel

	for i := 0; i < len(careers); i++ {
		tags, err = gs.GeneralRepo.GetTagsForUserOrCareer(careers[i].ID, false)
		if err != nil {
			panic(exceptions.Exception{
				Tag:    enums.INTERNAL_ERROR,
				Errors: []enums.SpecificError{enums.DATABASE_ERROR},
			})
		}
		res.Tags = append(res.Tags, tags)
	}
	return &res
}

func (gs *GeneralService) GetTagsForUser(userid int) []models.TagModel {
	res, err := gs.GeneralRepo.GetTagsForUserOrCareer(userid, true)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
		})
	}

	return res
}

func (gs *GeneralService) AddCareer(userid int, params *dto.PostCareerDTO) *models.CareerModel {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
		})
	}

	_, err := gs.UserRepo.FindUserByID(userid)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.NOT_FOUND,
			Errors: []enums.SpecificError{enums.USER_NOT_FOUND},
		})
	}

}
