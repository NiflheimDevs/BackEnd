package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
)

type TagService interface {
	GetTags() []byte
	GetTagsForUser(userid int) []dto.GetTagDto
	UpdateTagsForCareerOrUser(careerUserid int, newTags []dto.RecieveTagDTO, isForUser bool) []dto.RecieveTagDTO
	UpdateTagsForProject(projectid int, newTags []int) []int
}
