package repositories

import (
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type LabelRepo interface {
	GetLabels() []models.LabelModel
	GetLabelInfo(labelID int) *models.LabelModel
}
