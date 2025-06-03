package cdcdto

import (
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
)

// type UserDataDto struct {
// 	ID          int       `json:"id"`
// 	FirstName   string    `json:"firstname"`
// 	LastName    string    `json:"lastname"`
// 	Username    string    `json:"username"`
// 	Email       string    `json:"email"`
// 	Bio         string    `json:"bio"`
// 	Phone       string    `json:"phone"`
// 	CreatedTime time.Time `json:"created_time"`
// }

// type UserCapturer struct {
// 	Type   string       `json:"op"`
// 	After  *UserDataDto `json:"after"`
// 	Before *UserDataDto `json:"before"`
// }

type UserCapturer struct {
	Type   string                         `json:"op"`
	After  *elasticmodel.UserElasticModel `json:"after"`
	Before *elasticmodel.UserElasticModel `json:"before"`
}
