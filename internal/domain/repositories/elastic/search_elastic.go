package elastic

import elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"

type SearchRepo interface {
	SearchUsers(req *elasticmodel.SearchRequest) ([]map[string]any, error)
	SearchProjects(req *elasticmodel.SearchRequest) ([]map[string]any, error)
	SearchUsersProjectsTeams(req *elasticmodel.SearchRequest) ([]elasticmodel.InnerHits, error)
	SearchTeams(req *elasticmodel.SearchRequest) ([]map[string]any, error)
}
