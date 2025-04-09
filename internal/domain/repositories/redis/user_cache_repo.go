package redis

//requires redis version >= 6 otherwise most package calls will result in syntax error
import (
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type UserCache interface {
	FindByUsername(username string) (string, error)
	FindByPhone(phonenumber string) (string, error)
	PostUserCreds(session string, userdata *models.UserCacheData)
	FindBySession(session string) (string, error)
	ClearUserCreds(phonenumber string, username string)
	PostSessionOTP(session string, userdata *models.UserCacheData)
	PostSessionFlag(session string, userID string)
	GetSessionFlag(session string) (string, error)
	PostPhone(phone string, session string)
	DeleteSession(session string)
}
