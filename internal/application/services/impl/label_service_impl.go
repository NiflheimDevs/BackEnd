package servicesimpl

import (
	"encoding/json"

	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type LabelService struct {
	LableRepo repositories.LabelRepo
	UserRepo  repositories.UserRepo
}

func NewLabelService(
	labelRepo repositories.LabelRepo,
	userRepo repositories.UserRepo,
) *LabelService {
	return &LabelService{
		LableRepo: labelRepo,
		UserRepo:  userRepo,
	}
}

func (ls *LabelService) GetLabels() []byte {
	labels := ls.LableRepo.GetLabels()

	marshaled, _ := json.Marshal(labels)

	return marshaled
}

func (ls *LabelService) GetLabelInfo(labelID int) *models.LabelModel {
	label := ls.LableRepo.GetLabelInfo(labelID)

	return label
}
