package elastic

import elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"

type TeamElastic interface {
	UpsertTeamDoc(teamem *elasticmodel.TeamElasticModel) error
	GetTeamDoc(teamid int) *elasticmodel.TeamElasticModel
	DeleteTeamDoc(teamid int) error
}
