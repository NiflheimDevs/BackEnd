package elasticmodel

type SearchResult struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Title      string            `json:"title"`
	Summary    string            `json:"summary"`
	Highlights map[string]string `json:"highlights,omitempty"`
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

var esResp struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Index  string                 `json:"_index"`
			ID     string                 `json:"_id"`
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
