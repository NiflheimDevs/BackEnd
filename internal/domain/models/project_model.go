package models

type ProjectModel struct {
	ID          int
	OwnerID     int
	Title       string
	Description string
	Label       int
	SelectedBid int
	Tags        []TagModel
	Duration    string
}
