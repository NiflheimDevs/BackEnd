package dto

import "time"

type LoginDTO struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileDTO struct {
	FirstName          string    `json:"firstname"`
	LastName           string    `json:"lastname"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Is_verified        bool      `json:"is_verified"`
	Bio                string    `json:"bio"`
	Phone              string    `json:"phonenumber"`
	HighProfilePicture string    `json:"high_profile"`
	LowProfilePicture  string    `json:"low_profile"`
	CreatedAt          time.Time `json:"created_at"`
	Rating             float64   `json:"rating"`
	CommentCount       int       `json:"comments"`
}

type UpdateUserDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Bio       string `json:"bio"`
}
