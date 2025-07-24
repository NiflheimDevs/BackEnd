package enums

type UrlTokenPurpose string

const (
	TeamInvite        UrlTokenPurpose = "team_invite"
	EmailVerification UrlTokenPurpose = "email_verification"
	PasswordReset     UrlTokenPurpose = "reset_password"
)
