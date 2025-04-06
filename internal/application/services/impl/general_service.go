package services

import (
	"encoding/json"

	"github.com/niflheimdevs/backend/internal/dto"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/models"
	"github.com/niflheimdevs/backend/internal/repositories"
	"github.com/niflheimdevs/backend/internal/utils"
)

// ! DUPLICATE CODE
// ! UTILS

type GeneralService struct {
	GeneralRepo *repositories.GeneralRepo
	UserRepo    *repositories.UserRepo
	// Utils       *utils.Utils
}

func NewGeneralService(
	GeneralRepo *repositories.GeneralRepo,
	userRepo *repositories.UserRepo,
	// utils *utils.Utils,
) *GeneralService {
	return &GeneralService{
		GeneralRepo: GeneralRepo,
		UserRepo:    userRepo,
		// Utils:       utils,
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

func (service *GeneralService) GetLabels() []byte {
	labels := service.GeneralRepo.GetLabels()

	marshaled, _ := json.Marshal(labels)

	return marshaled
}

func (service *GeneralService) GetLabelInfo(labelID int) *models.LabelModel {
	label := service.GeneralRepo.GetLabelInfo(labelID)

	return label
}

func (gs *GeneralService) GetCareerForUser(userid int, target int) []dto.SendCareerDTO {
	var res []dto.SendCareerDTO
	careers, err := gs.GeneralRepo.GetCareersForUser(target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
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
		temp.Tags, err = gs.GeneralRepo.GetTagsForUserOrCareer(careers[i].ID, false)
		if err != nil {
			panic(exceptions.Exception{
				Tag:    enums.INTERNAL_ERROR,
				Errors: []enums.SpecificError{enums.DATABASE_ERROR},
			})
		}

		res = append(res, temp)
	}

	return res
}

func (gs *GeneralService) GetTagsForUser(userid int) []dto.GetTagDto {
	res, err := gs.GeneralRepo.GetTagsForUserOrCareer(userid, true)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
		})
	}

	return res
}

func (gs *GeneralService) UpdateCareer(userid int, params []dto.CareerDTO) []dto.CareerDTO {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
		})
	}

	careers, err := gs.GeneralRepo.GetCareersForUser(userid)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
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
			params[i].ID, err = gs.GeneralRepo.CreateCareer(userid, &params[i])
			// couldn't create. remove from response
			if err != nil {
				newCareerSet[params[i].ID] = false
				params = utils.RemoveUnordered(params, &i)
			} else {
				params[i].Tags = gs.UpdateTagsForCareerOrUser(params[i].ID, params[i].Tags, false)
			}
		} else {
			if existingCareerSet[params[i].ID] != nil {
				if !params[i].IsEqualToModel(existingCareerSet[params[i].ID]) {
					gs.GeneralRepo.UpdateCareer(userid, &params[i])
				}
				gs.UpdateTagsForCareerOrUser(params[i].ID, params[i].Tags, false)
			} else {
				// ? this means that a career has id but it shouldn't !!
				newCareerSet[params[i].ID] = false
				params = utils.RemoveUnordered(params, &i)
			}
		}
	}

	for careerid, _ := range existingCareerSet {
		if !newCareerSet[careerid] {
			gs.GeneralRepo.DeleteCareer(userid, careerid)
		}
	}

	return params
}

func (gs *GeneralService) UpdateTagsForCareerOrUser(careerUserid int, newTags []dto.RecieveTagDTO, isForUser bool) []dto.RecieveTagDTO {

	if careerUserid < 0 {
		panic(
			exceptions.Exception{
				Tag: enums.UNPROCESSABLE,
			})
	}

	tags, err := gs.GeneralRepo.GetTagsForUserOrCareer(careerUserid, isForUser)

	if err != nil {
		return []dto.RecieveTagDTO{}
	}

	existingTagSet := make(map[int]*dto.GetTagDto)
	newTagSet := make(map[int]bool)
	//create set
	for i := 0; i < len(tags); i++ {
		existingTagSet[tags[i].ID] = &tags[i]
	}

	for i := 0; i < len(newTags); i++ {
		newTagSet[newTags[i].ID] = true
		if existingTagSet[newTags[i].ID] != nil {
			if !newTags[i].IsEqualToGetTagDTO(existingTagSet[newTags[i].ID]) {
				gs.GeneralRepo.UpdateTagForUserOrCareer(&newTags[i], careerUserid, isForUser)
			}
		} else {
			err = gs.GeneralRepo.AddTagToUserOrCareer(&newTags[i], careerUserid, isForUser)
			if err != nil {
				utils.RemoveUnordered(tags, &i)
				newTagSet[newTags[i].ID] = false
			}
		}
	}
	for tagid, _ := range existingTagSet {
		if !newTagSet[tagid] {
			gs.GeneralRepo.DeleteTagForCareerOrUserByID(tagid, careerUserid, isForUser)
		}
	}

	return newTags
}
