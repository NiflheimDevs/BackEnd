package elasticimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/elastic/go-elasticsearch/v9"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
	"github.com/niflheimdevs/backend/internal/utils"
)

type SearchElastic struct {
	Es *elasticsearch.Client
}

func NewSearchElastic(es *elasticsearch.Client) *SearchElastic {
	return &SearchElastic{
		Es: es,
	}
}

func (se *SearchElastic) SearchUserProjectTeam(req *elasticmodel.SearchRequest) ([]map[string]any, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var indices []string
	if len(req.Types) == 0 || utils.Contains(req.Types, "users") && utils.Contains(req.Types, "projects") {
		indices = []string{"users", "projects"}
	} else {
		for _, t := range req.Types {
			if t == "users" || t == "projects" {
				indices = append(indices, t)
			}
		}
	}

	from := (req.Page - 1) * req.Limit
	if from < 0 {
		from = 0
	}

	mustTagFilters := make([]any, 0, len(req.Tags))
	for _, role := range req.Tags {
		mustTagFilters = append(mustTagFilters, map[string]any{
			"term": map[string]any{"tags": role},
		})
	}

	query := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query": req.Query,
							"fields": []string{"users.username^3",
								"users.firstname^3",
								"users.lastname^3",
								"projects.title^3",
								"projects.description",
								"teams.title^3",
								"teams.description"},
							"fuzziness": "AUTO",
						},
					},
				},
				"filter": []any{
					map[string]any{
						"bool": map[string]any{
							"must": mustTagFilters,
						},
					},
				},
			},
		},
		"from": from,
		"size": req.Limit,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := se.Es.Search(
		se.Es.Search.WithContext(ctx),
		se.Es.Search.WithIndex(indices...),
		se.Es.Search.WithBody(&buf),
		se.Es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		return nil, errors.New(string(data))
	}

}
