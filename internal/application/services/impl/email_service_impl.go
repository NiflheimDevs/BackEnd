package servicesimpl

import (
	"errors"
	"log"
	"path/filepath"
	"strconv"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/pkg"
)

type EmailService struct {
	Env             *bootstrap.Env
	UserRepo        repositories.UserRepo
	TeamRepo        repositories.TeamRepo
	UrlTokenService services.UrlTokenService
}

func NewEmailService(
	userRepo repositories.UserRepo,
	env *bootstrap.Env,
	tr repositories.TeamRepo,
	uts services.UrlTokenService,
) *EmailService {
	return &EmailService{
		UserRepo:        userRepo,
		Env:             env,
		TeamRepo:        tr,
		UrlTokenService: uts,
	}
}

func (es *EmailService) getTemplatePath(filename string) string {
	return filepath.Join("internal", "application", "email_templates", filename)
}

// i made a decision to use the email service for sending team invites. pro: 1- it is handled in one thread 2- the code remains the same when i microservice it later \t con:depends on other services highly coupled
func (es *EmailService) SendTeamInvites(userid int, newMembers []int, teamid int64) error {
	if userid < 0 {
		log.Println("SendTeamInviteError: invalid user id.")
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_ACCESS_DENIED},
		})
		// return errors.New("invalid user id")
	}

	inviter, err := es.UserRepo.FindUserByID(userid)
	if err != nil {
		log.Println("SendTeamInviteError: inviter error. details:", err)
		return err
	}
	if inviter == nil {
		log.Println("SendTeamInviteError: inviter not found. userid:", userid)
		panic(exceptions.Exception{
			Tag:    exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{exceptions.USER_NOT_FOUND},
		})
		// return errors.New("inviter not found")
	}

	teamInfo := es.TeamRepo.GetTeam(teamid)
	if teamInfo == nil {
		log.Println("SendTeamInviteError: team error.")
		return errors.New("no team found")
	}
	for _, member := range newMembers {
		userInfo, err := es.UserRepo.FindUserByID(member)
		if err != nil {
			log.Panicln("SendTeamInviteError: new member error. details:", err)
			continue
		}
		if userInfo == nil {
			log.Println("SendTeamInviteError: new member not found. userid:", member)
			continue
		}
		if !userInfo.Is_verified {
			continue
		}
		token, err := es.UrlTokenService.GenerateTeamInviteToken(userInfo.ID, userid, teamid)
		if err != nil {
			log.Println("SendTeamInviteError: token generation error. details:", err)
			continue
		}

		data := dto.TeamInviteHTML{
			InviterUsername: inviter.Username,
			TeamName:        teamInfo.Title,
			RedirectURL:     "https://bidlancer.ir/team/invite/?teamid=" + strconv.FormatInt(teamid, 10) + "&userid=" + strconv.Itoa(userInfo.ID) + "&token=" + token,
		}
		err = es.sendTeamInviteEmail(userInfo.Email, &data)
		if err != nil {
			log.Println("SendTeamInviteError: email sending error. details:", err, "\nto:", userInfo.Email, "\ndata:", data)
			continue
		}
	}

	return nil
}

func (es *EmailService) sendTeamInviteEmail(to string, data *dto.TeamInviteHTML) error {
	templatePath := es.getTemplatePath("team_invite.html")
	htmlBody, err := pkg.RenderTemplate(templatePath, data)
	if err != nil {
		return err
	}
	from := es.Env.Email.MainAddr
	host := es.Env.Email.Host
	port, _ := strconv.Atoi(es.Env.Email.Port)
	username := es.Env.Email.Username
	password := es.Env.Email.Password
	msg := pkg.NewHTMLMessage(from, to, "A Team Needs You!", htmlBody)
	return pkg.SendEmail(host, port, username, password, msg)
}

func (es *EmailService) SendEmailVerificationEmail(userid int) error {
	user, err := es.UserRepo.FindUserByID(userid)
	if err != nil {
		log.Panicln("SendEmailVerificationEmail: user error. details:", err)
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	token, err := es.UrlTokenService.GenerateEmailVerificationToken(user.ID)
	if err != nil {
		log.Panicln("SendEmailVerificationEmail: token generation error. details:", err)
		return err
	}

	data := dto.EmailVerificationHTML{
		Username:        user.Username,
		VerificationURL: "https://bidlancer.ir/verify/email/?userid=" + strconv.Itoa(user.ID) + "&token=" + token,
	}
	err = es.sendEmailVerificationEmail(user.Email, &data)
	if err != nil {
		log.Println("SendEmailVerificationEmail: email sending error. details:", err, "\nto:", user.Email, "\ndata:", data)
		return err
	}
	return nil
}

func (es *EmailService) sendEmailVerificationEmail(to string, data *dto.EmailVerificationHTML) error {
	templatePath := es.getTemplatePath("email_verification.html")
	htmlBody, err := pkg.RenderTemplate(templatePath, data)
	if err != nil {
		log.Println("SendEmailVerificationEmail: template rendering error. details:", err)
		return err
	}
	from := es.Env.Email.MainAddr
	host := es.Env.Email.Host
	port, _ := strconv.Atoi(es.Env.Email.Port)
	username := es.Env.Email.Username
	password := es.Env.Email.Password
	msg := pkg.NewHTMLMessage(from, to, "Are you you?", htmlBody)
	return pkg.SendEmail(host, port, username, password, msg)
}

func (es *EmailService) SendAcceptedInvitationEmail(inviteeid, inviterid int, teamid int64) error {
	invitee, err := es.UserRepo.FindUserByID(inviteeid)
	if err != nil {
		log.Panicln("SendAcceptedInvitationEmail: invitee error. details:", err)
		return err
	}
	if invitee == nil {
		return errors.New("invitee not found")
	}

	inviter, err := es.UserRepo.FindUserByID(inviterid)
	if err != nil {
		log.Panicln("SendAcceptedInvitationEmail: inviter error. details:", err)
		return err
	}
	if inviter == nil {
		return errors.New("inviter not found")
	}

	teamInfo := es.TeamRepo.GetTeam(teamid)
	if teamInfo == nil {
		log.Println("SendAcceptedInvitationEmail: team error.")
		return errors.New("no team found")
	}
	data := dto.InvitationAcceptedHTML{
		InviterUsername: inviter.Username,
		InviteeUsername: invitee.Username,
		TeamName:        teamInfo.Title,
	}
	return es.sendAcceptedInvitationEmail(invitee.Email, &data)
}

func (es *EmailService) sendAcceptedInvitationEmail(to string, data *dto.InvitationAcceptedHTML) error {
	templatePath := es.getTemplatePath("invitation_accepted.html")
	htmlBody, err := pkg.RenderTemplate(templatePath, data)
	if err != nil {
		return err
	}
	from := es.Env.Email.MainAddr
	host := es.Env.Email.Host
	port, _ := strconv.Atoi(es.Env.Email.Port)
	username := es.Env.Email.Username
	password := es.Env.Email.Password
	msg := pkg.NewHTMLMessage(from, to, "They Said Yes!", htmlBody)
	return pkg.SendEmail(host, port, username, password, msg)
}

func (es *EmailService) SendNotification(userID int, subject string, body string) error {
	user, _ := es.UserRepo.FindUserByID(userID)

	if user.Is_verified {
		from := es.Env.Email.NotifAddr
		host := es.Env.Email.Host
		port, _ := strconv.Atoi(es.Env.Email.Port)
		username := es.Env.Email.Username
		password := es.Env.Email.Password
		message := pkg.NewMessage(from, user.Email, subject, body)
		err := pkg.SendEmail(host, port, username, password, message)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (es *EmailService) SendEmail(to, subject string, body string) error {
	from := es.Env.Email.MainAddr
	host := es.Env.Email.Host
	port, _ := strconv.Atoi(es.Env.Email.Port)
	username := es.Env.Email.Username
	password := es.Env.Email.Password
	message := pkg.NewMessage(from, to, subject, body)
	err := pkg.SendEmail(host, port, username, password, message)
	if err != nil {
		return err
	}
	return nil
}
