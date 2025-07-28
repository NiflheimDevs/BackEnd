package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
)

type UserService interface {
	CheckAvailabilityForSignup(phonenumber string, username string) string
	CacheUserInfo(phonenumber string, username string, password string, otp string) string
	ValidateOTP(session string, otp string) (string, string, []byte)
	Register(phonenumber string, username string, password []byte, session string) int
	AuthenticateUser(identifier string, password string) *models.UserModel
	ChangePasswordValidate(user_id int, old_password string)
	ChangePassword(user_id int, new_password string)
	SetupOTP(phonenumber string, code string) string
	SetForgetPasswordFlag(userID string) string
	CheckFlagForPasswordReset(session string) int
	UpdateUserData(userid int, userData *dto.UpdateUserDTO)
	UpdateEmail(userid int, email string)
	UpdateUsername(userid int, username string)
	UpdatePhoneSendOTP(phone string, userid int, code string) string
	UpdatePhone(phone string, userid string)
	GetUserInfo(targetUserid int, userid int) *dto.UserProfileDTO
	DeleteUser(userid int)
	SearchUsers(req *elasticmodel.QueryAndTagSearchReqDto) []map[string]any

	VerifyEmail(userid int, token string)
	ResendEmailVerification(userid int)
}
