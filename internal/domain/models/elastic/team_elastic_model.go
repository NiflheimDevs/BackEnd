package elasticmodel

import "time"

type TeamElasticModel struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Type        int       `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
	Created_at  time.Time `json:"created_at,omitempty"`
}
