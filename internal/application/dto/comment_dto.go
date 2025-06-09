package dto

type CommentDTO struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Content   string `json:"content"`
	Rating    int    `json:"rating"`
	UserID    int    `json:"user_id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
}

type PutComment struct {
	ID int `json:"id"`
}
