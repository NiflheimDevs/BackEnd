package elasticimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
)

type SearchRepo struct {
	Es *elasticsearch.Client
}

func NewSearchElastic(es *elasticsearch.Client) *SearchRepo {
	return &SearchRepo{
		Es: es,
	}
}

func (se *SearchRepo) SearchUsers(req *elasticmodel.SearchRequest) ([]map[string]any, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	from := (req.Page - 1) * req.Limit
	if from < 0 {
		from = 0
	}

	mustTagFilters := make([]any, 0, len(req.Tags))
	for _, tag := range req.Tags {
		mustTagFilters = append(mustTagFilters, map[string]any{
			"term": map[string]any{"tags.name": tag},
		})
	}

	query := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query": req.Query,
							"fields": []string{
								"username^3",
								"firstname^3",
								"lastname^3",
							},
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
		"sort": []any{
			map[string]any{
				req.SortBy: map[string]any{
					"order": req.Order,
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
		se.Es.Search.WithIndex(enums.ELASTIC_USER_INDEX),
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

	var respond elasticmodel.SearchResult

	if err = json.NewDecoder(res.Body).Decode(&respond); err != nil {
		log.Println("CastError: failed to cast response to struct. details:", err)
		return nil, err
	}

	return se.extractInnerHits(&respond), nil
}

func (se *SearchRepo) SearchProjects(req *elasticmodel.SearchRequest) ([]map[string]any, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	from := (req.Page - 1) * req.Limit
	if from < 0 {
		from = 0
	}

	mustTagFilters := make([]any, 0, len(req.Tags))
	for _, tag := range req.Tags {
		mustTagFilters = append(mustTagFilters, map[string]any{
			"term": map[string]any{"tags.name": tag},
		})
	}

	query := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query": req.Query,
							"fields": []string{
								"title^3",
								"description",
							},
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
		"sort": []any{
			map[string]any{
				req.SortBy: map[string]any{
					"order": req.Order,
				},
			},
			map[string]any{
				"projects.created_time": map[string]any{
					"order": "asc",
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
		se.Es.Search.WithIndex(enums.ELASTIC_PROJECT_INDEX),
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

	var respond elasticmodel.SearchResult

	if err = json.NewDecoder(res.Body).Decode(&respond); err != nil {
		log.Println("CastError: failed to cast response to struct. details:", err)
		return nil, err
	}

	return se.extractInnerHits(&respond), nil
}

func (se *SearchRepo) SearchUsersProjectsTeams(req *elasticmodel.SearchRequest) ([]elasticmodel.InnerHits, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var indices []string

	for _, t := range req.Types {
		if t == "users" {
			indices = append(indices, enums.ELASTIC_USER_INDEX)
		} else if t == "projects" {
			indices = append(indices, enums.ELASTIC_PROJECT_INDEX)
		} else if t == "teams" {
			indices = append(indices, enums.ELASTIC_TEAM_INDEX)
		}
	}

	from := (req.Page - 1) * req.Limit
	if from < 0 {
		from = 0
	}

	query := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query": req.Query,
				"fields": []string{
					"username^3",
					"firstname^3",
					"lastname^3",
					"title^3",
					"description",
				},
				"fuzziness": "AUTO",
			},
		},
		"sort": []any{
			map[string]any{
				"_score": map[string]any{
					"order": "desc",
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

	var respond elasticmodel.SearchResult

	if err = json.NewDecoder(res.Body).Decode(&respond); err != nil {
		log.Println("CastError: failed to cast response to struct. details:", err)
		return nil, err
	}
	return respond.Hits.Hits, nil
}

func (se *SearchRepo) SearchTeams(req *elasticmodel.SearchRequest) ([]map[string]any, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	from := (req.Page - 1) * req.Limit
	if from < 0 {
		from = 0
	}

	query := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query": req.Query,
				"fields": []string{
					"title^3",
					"description"},
				"fuzziness": "AUTO",
			},
		},
		"sort": []any{
			map[string]any{
				req.SortBy: map[string]any{
					"order": req.Order,
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
		se.Es.Search.WithIndex(enums.ELASTIC_TEAM_INDEX),
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

	var respond elasticmodel.SearchResult

	if err = json.NewDecoder(res.Body).Decode(&respond); err != nil {
		log.Println("CastError: failed to cast response to struct. details:", err)
		return nil, err
	}
	return se.extractInnerHits(&respond), nil
}

func (se *SearchRepo) extractInnerHits(respond *elasticmodel.SearchResult) []map[string]any {

	var result []map[string]any

	for i := 0; i < respond.Hits.Total.Value; i++ {
		result = append(result, respond.Hits.Hits[i].Source)
	}

	return result

}
