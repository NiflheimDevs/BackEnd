package models

type CommentModel struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	BidID     int    `json:"bid_id"`
	Content   string `json:"content"`
	Rating    int    `json:"rating"`
}
