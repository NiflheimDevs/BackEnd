package repositories

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type CareerRepo interface {
	GetCareersForUser(userid int) ([]models.CareerModel, error)
	CreateCareer(userid int, params *dto.CareerDTO) (int, error)
	UpdateCareer(userid int, params *dto.CareerDTO) error
	DeleteCareer(userid int, careerid int) error
}
