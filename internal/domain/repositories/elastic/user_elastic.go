package elastic

import elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"

type UserElastic interface {
	UpsertUserDoc(userem *elasticmodel.UserElasticModel) error
	GetUserDoc(userid int) *elasticmodel.UserElasticModel
	DeleteUserDoc(userid int) error
}
