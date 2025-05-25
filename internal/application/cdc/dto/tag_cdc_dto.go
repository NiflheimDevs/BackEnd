package cdcdto

type ProjectTagDto struct {
	Projectid int `json:"project_id"`
	Tagid     int `json:"tag_id"`
}

type ProjectTagCapturer struct {
	Type   string         `json:"op"`
	After  *ProjectTagDto `json:"after"`
	Before *ProjectTagDto `json:"before"`
}

type UserCareerTagDto struct {
	UserCareerid int `json:"career_user_id"`
	Tagid        int `json:"tag_id"`
	Type         int `json:"type"`
	Level        int `json:"level"`
}

type UserCareerTagCapturer struct {
	Type   string            `json:"op"`
	After  *UserCareerTagDto `json:"after"`
	Before *UserCareerTagDto `json:"before"`
}
