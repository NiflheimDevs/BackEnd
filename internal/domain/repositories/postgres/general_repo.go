package repositories

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type GeneralRepo interface {
	GetTags() ([]models.TagModel, error)

	GetLabels() []models.LabelModel

	GetLabelInfo(labelID int) *models.LabelModel

	GetTagsForUserOrCareer(careerUserid int, isForUser bool) ([]dto.GetTagDto, error)

	DeleteTagForUserOrCareer(tagid int, careerUserid int, isForUser bool) error

	UpdateTagForUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error

	AddTagToUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error

	GetCareersForUser(userid int) ([]models.CareerModel, error)

	CreateCareer(userid int, params *dto.CareerDTO) (int, error)

	UpdateCareer(userid int, params *dto.CareerDTO) error

	DeleteTagForCareerOrUserByID(tagid int, careerUserid int, isForUser bool) error

	DeleteCareer(userid int, careerid int) error
}
