package servicesimpl

import (
	"log"

	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
	"github.com/niflheimdevs/backend/internal/domain/repositories/elastic"
)

type SherlockService struct {
	SearchRepo  elastic.SearchRepo
	FileService services.FileService
}

func NewSherlockService(
	sr elastic.SearchRepo,
	fs services.FileService,
) *SherlockService {
	return &SherlockService{
		SearchRepo:  sr,
		FileService: fs,
	}
}

func (ss *SherlockService) SearchEverything(req *elasticmodel.SearchRequest) []elasticmodel.InnerHits {

	if req.SortBy == "" {
		req.SortBy = "_score"
	}
	if req.Order == "" {
		req.Order = "desc"
	}
	if len(req.Types) == 0 || (len(req.Tags) == 1 && req.Tags[0] == "") {
		req.Types = []string{enums.ELASTIC_USER_INDEX, enums.ELASTIC_PROJECT_INDEX, enums.ELASTIC_TEAM_INDEX}
	}

	res, err := ss.SearchRepo.SearchUsersProjectsTeams(req)
	if err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

	for i := 0; i < len(res); i++ {
		if res[i].Index == enums.ELASTIC_USER_INDEX {
			res[i].Source["profile"] = ss.FileService.GetProfilePhotoURL(int(res[i].Source["id"].(float64)), false)
		} else if res[i].Index == enums.ELASTIC_TEAM_INDEX {
			res[i].Source["profile"] = ss.FileService.GetTeamProfilePhotoURL(int64(res[i].Source["id"].(float64)), false)
		}
	}

	return res
}
