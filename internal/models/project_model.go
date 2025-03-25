package models

type ProjectModel struct {
	ID          int
	OwnerID     int
	Title       string
	Description string
	Label       []string
	Tags        []TagModel
	Duration    string
}
