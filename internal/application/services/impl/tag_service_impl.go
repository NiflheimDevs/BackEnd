package servicesimpl

import (
	"encoding/json"
	"log"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/utils"
)

type TagService struct {
	TagRepo  repositories.TagRepo
	UserRepo repositories.UserRepo
}

func NewTagService(
	tagRepo repositories.TagRepo,
	userRepo repositories.UserRepo,
) *TagService {
	return &TagService{
		TagRepo:  tagRepo,
		UserRepo: userRepo,
	}
}

func (ts *TagService) GetTags() []byte {
	tags, err := ts.TagRepo.GetTags()
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	marshaled, err := json.Marshal(tags)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}

	return marshaled
}

func (ts *TagService) GetTagsForUser(userid int) []dto.GetTagDto {
	res, err := ts.TagRepo.GetTagsForUserOrCareer(userid, true)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	return res
}

func (ts *TagService) UpdateTagsForCareerOrUser(careerUserid int, newTags []dto.RecieveTagDTO, isForUser bool) []dto.RecieveTagDTO {

	if careerUserid < 0 {
		panic(
			exceptions.Exception{
				Tag: exceptions.UNPROCESSABLE,
			})
	}

	tags, err := ts.TagRepo.GetTagsForUserOrCareer(careerUserid, isForUser)

	if err != nil {
		return []dto.RecieveTagDTO{}
	}

	existingTagSet := make(map[int]*dto.GetTagDto)
	newTagSet := make(map[int]bool)
	//create set
	for i := 0; i < len(tags); i++ {
		existingTagSet[tags[i].ID] = &tags[i]
	}

	log.Println("existing tags: ", existingTagSet)

	for i := 0; i < len(newTags); i++ {
		newTagSet[newTags[i].ID] = true
		if existingTagSet[newTags[i].ID] != nil {
			log.Println(newTagSet, "exists")
			if !newTags[i].IsEqualToGetTagDTO(existingTagSet[newTags[i].ID]) {
				log.Println("updating")
				err = ts.TagRepo.UpdateTagForUserOrCareer(&newTags[i], careerUserid, isForUser)
				if err != nil {
					log.Println("Error Updating tag", newTags[i], "details :", err)
				}
			}
		} else {
			err = ts.TagRepo.AddTagToUserOrCareer(&newTags[i], careerUserid, isForUser)
			if err != nil {
				utils.RemoveUnordered(tags, &i)
				newTagSet[newTags[i].ID] = false
				log.Println("failed to add tag", newTags[i], "details:", err)
			}
		}
	}
	for tagid, _ := range existingTagSet {
		if !newTagSet[tagid] {
			err = ts.TagRepo.DeleteTagForCareerOrUserByID(tagid, careerUserid, isForUser)
			if err != nil {
				log.Println("failed to delete tag", tagid, "details:", err)
			}
		}
	}

	return newTags
}

func (ts *TagService) UpdateTagsForProject(projectid int, newTags []int) {

	existingTags := ts.TagRepo.GetProjectTag(projectid)

	existingTagSet := make(map[int]bool)
	for _, tag := range existingTags {
		existingTagSet[tag.ID] = true
	}

	newTagSet := make(map[int]bool)
	for _, tag := range newTags {
		newTagSet[tag] = true
	}

	for _, tag := range newTags {
		if !existingTagSet[tag] {
			ts.TagRepo.AddProjectTag(projectid, tag)
		}
	}

	for _, tag := range existingTags {
		if !newTagSet[tag.ID] {
			ts.TagRepo.DeleteProjectTag(projectid, tag.ID)
		}
	}
}
