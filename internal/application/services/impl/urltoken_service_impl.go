package servicesimpl

import (
	"errors"
	"time"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/utils"
)

type UrlTokenService struct {
	Constants    *bootstrap.Constants
	UrlTokenRepo repositories.UrlTokenRepo
}

func NewUrlTokenService(r repositories.UrlTokenRepo, c *bootstrap.Constants) *UrlTokenService {
	return &UrlTokenService{
		UrlTokenRepo: r,
		Constants:    c,
	}
}

func (s *UrlTokenService) GenerateTeamInviteToken(userid, senderID int, teamid int64) (string, error) { // , invitedEmail string

	if senderID < 0 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_TOKEN_EXPIRED},
		})
	}

	s.UrlTokenRepo.DeleteTokenByTeamAndPurposeAndUser(teamid, userid, enums.TeamInvite)

	tokenStr, err := utils.GenerateToken(s.Constants.UrlTokenSetting.Length)
	if err != nil {
		return "", err
	}

	token := &models.UrlToken{
		Token:   tokenStr,
		Purpose: enums.TeamInvite,
		// InvitedEmail: invitedEmail,
		TeamID:    &teamid,
		SenderID:  &senderID,
		UserID:    &userid,
		ExpiresAt: time.Now().Add(s.Constants.UrlTokenSetting.ExpiryTeamInvite),
	}
	if err := s.UrlTokenRepo.InsertToken(token); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *UrlTokenService) GenerateEmailVerificationToken(userid int) (string, error) {

	if userid < 0 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_TOKEN_EXPIRED},
		})
	}

	s.UrlTokenRepo.DeleteTokenByPurposeAndUser(userid, enums.EmailVerification)
	tokenStr, err := utils.GenerateToken(s.Constants.UrlTokenSetting.Length)
	if err != nil {
		return "", err
	}

	token := &models.UrlToken{
		Token:     tokenStr,
		Purpose:   enums.EmailVerification,
		UserID:    &userid,
		ExpiresAt: time.Now().Add(s.Constants.UrlTokenSetting.ExpirtyEmailVer),
	}
	if err := s.UrlTokenRepo.InsertToken(token); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *UrlTokenService) GetToken(tokenStr string) (*models.UrlToken, error) { //, expectedPurpose enums.UrlTokenPurpose
	urltoken, err := s.UrlTokenRepo.GetTokenByValue(tokenStr)
	if err != nil || urltoken == nil {
		return nil, errors.New("token not found")
	}

	if urltoken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return urltoken, nil
}

func (s *UrlTokenService) GetTokenWithPurpose(tokenStr string, expectedPurpose enums.UrlTokenPurpose) (*models.UrlToken, error) {
	urltoken, err := s.UrlTokenRepo.GetTokenByValueAndPurpose(tokenStr, expectedPurpose)
	if err != nil || urltoken == nil {
		return nil, errors.New("token not found")
	}

	if urltoken.ExpiresAt.Before(time.Now()) {
		return urltoken, errors.New("token expired")
	}
	return urltoken, nil
}

func (s *UrlTokenService) DeleteToken(token string) error {
	err := s.UrlTokenRepo.DeleteToken(token)
	if err != nil {
		return err
	}
	return nil
}
