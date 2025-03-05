package dto

type LoginDTO struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
