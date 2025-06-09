package elasticmodel

type SearchResult struct {
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64     `json:"max_score"`
		Hits     []InnerHits `json:"hits"`
	} `json:"hits"`
}

type InnerHits struct {
	Index  string         `json:"_index"`
	ID     string         `json:"_id"`
	Source map[string]any `json:"_source"`
}

type SearchRequest struct {
	Query  string   `json:"query"`
	Tags   []string `json:"tags"`
	Types  []string `json:"types"`
	Page   int      `json:"page"`
	Limit  int      `json:"limit"`
	SortBy string   `json:"sort_by"`
	Order  string   `json:"order"`
}

type SimpleQuerySearchReqDto struct {
	Query  string `json:"query" validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Limit  int    `json:"limit" validate:"required"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}

type QueryAndTagSearchReqDto struct {
	Query  string   `json:"query" validate:"required"`
	Tags   []string `json:"tags"`
	Page   int      `json:"page" validate:"required"`
	Limit  int      `json:"limit" validate:"required"`
	SortBy string   `json:"sort_by"`
	Order  string   `json:"order"`
}
