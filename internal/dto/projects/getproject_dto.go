package dto

type Project struct {
	ProjectID   int    `json:"project_id"`
	OwnerID     int    `json:"Owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
