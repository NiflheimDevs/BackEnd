package dto

type TeamInviteHTML struct {
	InviterUsername string `json:"InviterUsername"`
	TeamName        string `json:"TeamName"`
	RedirectURL     string `json:"RedirectURL"`
}

type EmailVerificationHTML struct {
	Username        string `json:"Username"`
	VerificationURL string `json:"VerificationURL"`
}

type InvitationAcceptedHTML struct {
	InviterUsername string `json:"InviterUsername"`
	InviteeUsername string `json:"InviteeUsername"`
	TeamName        string `json:"TeamName"`
}
