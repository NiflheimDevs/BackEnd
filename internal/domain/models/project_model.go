package models

type ProjectModel struct {
	ID          int
	OwnerID     int
	Title       string
	Description string
	Label       int
	Tags        []TagModel
	Duration    string
}
