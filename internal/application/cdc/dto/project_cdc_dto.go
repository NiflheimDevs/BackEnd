package cdcdto

import elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"

// type ProjectDto struct {
// 	ID          int       `json:"id"`
// 	OwnerID     int       `json:"owner_id"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	Label       int       `json:"label"`
// 	State       int       `json:"state"`
// 	CreatedTime time.Time `json:"created_time"`
// 	UpdatedTime time.Time `json:"updated_time"`
// 	Duration    time.Time `json:"duration"`
// 	StartTime   time.Time `json:"start_time"`
// 	EndTime     time.Time `json:"end_time"`
// }

// type ProjectCapturer struct {
// 	Type   string      `json:"op"`
// 	After  *ProjectDto `json:"after"`
// 	Before *ProjectDto `json:"before"`
// }

type ProjectCapturer struct {
	Type   string                            `json:"op"`
	After  *elasticmodel.ProjectElasticModel `json:"after"`
	Before *elasticmodel.ProjectElasticModel `json:"before"`
}
