package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type TagRepo interface {
	GetTags() ([]models.TagModel, error)
	GetTagsForUserOrCareer(careerUserid int, isForUser bool) ([]dto.GetTagDto, error)
	DeleteTagForUserOrCareer(tagid int, careerUserid int, isForUser bool) error
	UpdateTagForUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error
	GetProjectTag(projectID int) []models.TagModel
	AddProjectTagWithTx(ctx context.Context, tx transaction.Tx, projectID, tagID int) error
	AddProjectTag(projectID, tagID int) error
	DeleteProjectTagWithTx(ctx context.Context, tx transaction.Tx, projectID, tagID int) error
	DeleteProjectTag(projectID, tagID int) error
	AddTagToUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error
	DeleteTagForCareerOrUserByID(tagid int, careerUserid int, isForUser bool) error
}
