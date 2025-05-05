package models

import "time"

type ProjectModel struct {
	ID          int
	OwnerID     int
	Title       string
	Description string
	Label       int
	SelectedBid int
	State       int
	Tags        []TagModel
	Duration    time.Time
}
