package elasticmodel

import (
	"time"
)

type ProjectElasticModel struct {
	ID          int       `json:"id,omitempty"`
	OwnerID     int       `json:"owner_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Label       int       `json:"label,omitempty"`
	State       int       `json:"state,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
	UpdatedTime time.Time `json:"updated_time,omitempty"`
	Duration    time.Time `json:"duration,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}
