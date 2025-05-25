package cdcdto

import (
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
)

// type TeamCdcMessage struct {
// 	Type   string            `json:"op"`
// 	After  *models.TeamModel `json:"after"`
// 	Before *models.TeamModel `json:"before"`
// }

type TeamCdcMessage struct {
	Type   string                         `json:"op"`
	After  *elasticmodel.TeamElasticModel `json:"after"`
	Before *elasticmodel.TeamElasticModel `json:"before"`
}
