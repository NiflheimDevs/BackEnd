package models

type CommentModel struct {
	ID        int
	ProjectID int
	BidID     int
	Content   string
	Rating    int
}
