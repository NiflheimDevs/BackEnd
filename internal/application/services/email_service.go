package services

type EmailService interface {
	SendTeamInvites(userId int, newMembers []int, teamId int64) error
	// SendEmailVerificationEmail(to string, data interface{}) error
	// SendTeamInviteEmail(to string, data interface{}) error

	SendAcceptedInvitationEmail(inviteeid, inviterid int, teamid int64) error
	SendEmailVerificationEmail(userid int) error

	SendNotification(userId int, subject string, body string) error
	SendEmail(to, subject string, body string) error
}
