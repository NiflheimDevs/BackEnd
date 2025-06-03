package elasticmodel

import (
	"time"
)

type UserElasticModel struct {
	ID          int               `json:"id,omitempty"`
	FirstName   string            `json:"firstname,omitempty"`
	LastName    string            `json:"lastname,omitempty"`
	Username    string            `json:"username,omitempty"`
	Email       string            `json:"email,omitempty"`
	Bio         string            `json:"bio,omitempty"`
	Phone       string            `json:"phone,omitempty"`
	CreatedTime time.Time         `json:"created_time,omitempty"`
	Tags        []TagElasticModel `json:"tags,omitempty"`
}
