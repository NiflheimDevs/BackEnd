package cdcdto

type ChangeDataCaptureDto[T any] struct {
	Type   string `json:"op"`
	After  *T     `json:"after"`
	Before *T     `json:"before"`
}

type ProjectTagDto struct {
	Projectid int `json:"project_id"`
	Tagid     int `json:"tag_id"`
}

type UserCareerTagDto struct {
	UserCareerid int `json:"career_user_id"`
	Tagid        int `json:"tag_id"`
	Type         int `json:"type"`
	Level        int `json:"level"`
}

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
