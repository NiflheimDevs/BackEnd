package models

type UserModel struct {
	ID          int
	FirstName   string
	LastName    string
	Username    string
	password    string
	email       string
	is_verified bool
	bio         string
	phone       string
}
