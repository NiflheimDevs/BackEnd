package services

import (
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type LabelService interface {
	GetLabels() []byte
	GetLabelInfo(labelID int) *models.LabelModel
}
