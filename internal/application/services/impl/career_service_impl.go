package servicesimpl

import (
	"log"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/utils"
)

type CareerService struct {
	TagService services.TagService
	CareerRepo repositories.CareerRepo
	TagRepo    repositories.TagRepo
	UserRepo   repositories.UserRepo
}

func NewCareerService(
	careerRepo repositories.CareerRepo,
	tagService services.TagService,
	tagRepo repositories.TagRepo,
	userRepo repositories.UserRepo,
) *CareerService {
	return &CareerService{
		TagService: tagService,
		UserRepo:   userRepo,
		CareerRepo: careerRepo,
		TagRepo:    tagRepo,
	}
}

func (cs *CareerService) GetCareersForUser(userid int, target int) []dto.SendCareerDTO {
	var res []dto.SendCareerDTO
	careers, err := cs.CareerRepo.GetCareersForUser(target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	var temp dto.SendCareerDTO

	for i := 0; i < len(careers); i++ {
		temp = dto.SendCareerDTO{
			ID:        -1,
			StartDate: careers[i].StartDate,
			EndDate:   careers[i].EndDate,
			Company:   careers[i].Company,
			Role:      careers[i].Role,
			Website:   careers[i].Website,
		}
		if target == userid {
			temp.ID = careers[i].ID
		}
		temp.Tags, err = cs.TagRepo.GetTagsForUserOrCareer(careers[i].ID, false)
		if err != nil {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}

		res = append(res, temp)
	}

	return res
}

func (cs *CareerService) UpdateCareers(userid int, params []dto.CareerDTO) []dto.CareerDTO {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
		})
	}

	careers, err := cs.CareerRepo.GetCareersForUser(userid)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	existingCareerSet := make(map[int]*models.CareerModel)
	//create set
	for i := 0; i < len(careers); i++ {
		existingCareerSet[careers[i].ID] = &careers[i]
	}

	newCareerSet := make(map[int]bool)

	for i := 0; i < len(params); i++ {
		// create check
		newCareerSet[params[i].ID] = true
		if params[i].ID == -1 {
			params[i].ID, err = cs.CareerRepo.CreateCareer(userid, &params[i])
			// couldn't create. remove from response
			if err != nil {
				log.Println("CareerError: can't create for user", userid, ". detail:", err)
				newCareerSet[params[i].ID] = false
				params = utils.RemoveUnordered(params, &i)
			} else {
				params[i].Tags = cs.TagService.UpdateTagsForCareerOrUser(params[i].ID, params[i].Tags, false)
			}
		} else {
			if existingCareerSet[params[i].ID] != nil {
				if !params[i].IsEqualToModel(existingCareerSet[params[i].ID]) {
					cs.CareerRepo.UpdateCareer(userid, &params[i])
				}
				cs.TagService.UpdateTagsForCareerOrUser(params[i].ID, params[i].Tags, false)
			} else {
				// ? this means that a career has id but it shouldn't !!
				newCareerSet[params[i].ID] = false
				params = utils.RemoveUnordered(params, &i)
			}
		}
	}

	for careerid, _ := range existingCareerSet {
		if !newCareerSet[careerid] {
			cs.CareerRepo.DeleteCareer(userid, careerid)
		}
	}

	return params
}
