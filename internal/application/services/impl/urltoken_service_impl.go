package servicesimpl

import (
	"context"
	"errors"
	"time"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/domain/enums"
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

func (s *UrlTokenService) GenerateTeamInviteToken(ctx context.Context, userid, teamID, senderID int, invitedEmail string) (string, error) {
	tokenStr, err := utils.GenerateToken(s.Constants.UrlTokenSetting.Length)
	if err != nil {
		return "", err
	}

	token := &models.UrlToken{
		Token:        tokenStr,
		Purpose:      enums.TeamInvite,
		InvitedEmail: invitedEmail,
		TeamID:       &teamID,
		SenderID:     &senderID,
		UserID:       &userid,
		ExpiresAt:    time.Now().Add(s.Constants.UrlTokenSetting.ExpiryTeamInvite),
	}
	if err := s.UrlTokenRepo.InsertToken(ctx, token); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *UrlTokenService) GenerateEmailVerificationToken(ctx context.Context, userid int) (string, error) {
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
	if err := s.UrlTokenRepo.InsertToken(ctx, token); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *UrlTokenService) GetToken(ctx context.Context, tokenStr string, expectedPurpose enums.UrlTokenPurpose) (*models.UrlToken, error) {
	urltoken, err := s.UrlTokenRepo.GetTokenByValue(ctx, tokenStr)
	if err != nil || urltoken == nil {
		return nil, errors.New("token not found")
	}

	if urltoken.Purpose != expectedPurpose {
		return nil, errors.New("invalid purpose")
	}

	if urltoken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired or already used")
	}

	return urltoken, nil
}
