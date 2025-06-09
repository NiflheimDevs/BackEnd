package services

import elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"

type SherlockService interface {
	SearchEverything(req *elasticmodel.SearchRequest) []elasticmodel.InnerHits
}
