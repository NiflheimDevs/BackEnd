package dto

import (
	"time"

	"github.com/niflheimdevs/backend/internal/models"
)

type SendCareerDTO struct {
	ID        int         `json:"id"`
	Company   string      `json:"company"`
	StartDate time.Time   `json:"start_date"`
	EndDate   time.Time   `json:"end_date"`
	Role      string      `json:"role"`
	Website   string      `json:"website"`
	Tags      []GetTagDto `json:"tags"`
}

type CareerDTO struct {
	ID        int             `json:"id" validate:"required,gt=-2"`
	Company   string          `json:"company" validate:"required,lt=25,gt=1"`
	StartDate time.Time       `json:"start_date" validate:"required"`
	EndDate   time.Time       `json:"end_date"`
	Role      string          `json:"role" validate:"required,lt=20,gt=1"`
	Website   string          `json:"website"`
	Tags      []RecieveTagDTO `json:"tags" validate:"dive"`
}

func (c *CareerDTO) IsEqualToModel(cm *models.CareerModel) bool {
	return c.Company == cm.Company &&
		c.EndDate == cm.EndDate &&
		c.Role == cm.Role &&
		c.StartDate == cm.StartDate &&
		c.Website == cm.Website
}
