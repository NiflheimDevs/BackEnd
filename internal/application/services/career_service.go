package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
)

type CareerService interface {
	GetCareersForUser(userid int, target int) []dto.SendCareerDTO
	UpdateCareers(userid int, params []dto.CareerDTO) []dto.CareerDTO
}
