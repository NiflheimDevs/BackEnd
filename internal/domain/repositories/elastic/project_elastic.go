package elastic

import elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"

type ProjectElastic interface {
	UpsertProjectDoc(projectem *elasticmodel.ProjectElasticModel) error
	GetProjectDoc(projectid int) *elasticmodel.ProjectElasticModel
	DeleteProjectDoc(projectid int) error
}
