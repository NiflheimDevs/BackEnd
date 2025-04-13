package dto

type LoginDTO struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileDTO struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Is_verified    bool   `json:"is_verified"`
	Bio            string `json:"bio"`
	Phone          string `json:"phonenumber"`
	ProfilePicture string `json:"profile"`
}

type UpdateUserDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Bio       string `json:"bio"`
}
