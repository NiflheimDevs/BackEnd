package dto

import (
	"time"

	"github.com/niflheimdevs/backend/internal/models"
)

type SendCareersDTO struct {
	Careers []models.CareerModel `json:"careers"`
	Tags    [][]GetTagDto        `json:"tags"`
}

type PostCareerDTO struct {
	Company   string    `json:"company" validate:"required,lt=25,gt=1"`
	StartDate time.Time `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   time.Time `json:"end_date"`
	Role      string    `json:"role" validate:"required,lt=20,gt=1"`
	Website   string    `json:"website"`
	Tags      []int     `json:"tags" validate:"dive,numeric"`
}
