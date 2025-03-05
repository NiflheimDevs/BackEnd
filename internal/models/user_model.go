package models

type UserModel struct {
	ID          int
	FirstName   string
	LastName    string
	Username    string
	Password    string
	Email       string
	Is_verified bool
	Bio         string
	Phone       string
}
