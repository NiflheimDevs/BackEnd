package elasticmodel

type ElasticsearchHit[T any] struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
}
