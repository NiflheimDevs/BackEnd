package redisimpl

//requires redis version >= 6 otherwise most package calls will result in syntax error
import (
	"context"
	"encoding/json"
	"time"

	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	DB *redis.Client
}

func NewUserCache(DB *redis.Client) *UserCache {
	return &UserCache{
		DB: DB,
	}
}

// ? should i make constant or functions for setting or giving the keys?

func PostRedis(uc *UserCache, key string, value interface{}, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	uc.DB.Set(ctx, key, value, duration)
}
func DeleteRedis(uc *UserCache, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uc.DB.Del(ctx, key)
}
func GetRedis(uc *UserCache, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return uc.DB.Get(ctx, key).Result()
}

func (uc *UserCache) FindByUsername(username string) (string, error) {

	return GetRedis(uc, "username:"+username)
}

func (uc *UserCache) FindByPhone(phonenumber string) (string, error) {

	return GetRedis(uc, "phone:"+phonenumber)
}

// ? idk if i should break this down
func (uc *UserCache) PostUserCreds(session string, userdata *models.UserCacheData) {

	PostRedis(uc, "username:"+userdata.Username, session, time.Minute*2+time.Second*2)
	PostRedis(uc, "phone:"+userdata.Phone, session, time.Minute*2+time.Second*2)

	val, _ := json.Marshal(*userdata)
	PostRedis(uc, session, val, time.Minute*2+time.Second*2)

}

func (uc *UserCache) FindBySession(session string) (string, error) {

	return GetRedis(uc, session)
}

func (uc *UserCache) ClearUserCreds(phonenumber string, username string) {
	DeleteRedis(uc, "username:"+username)

	DeleteRedis(uc, "phone:"+phonenumber)

}

func (uc *UserCache) PostSessionOTP(session string, userdata *models.UserCacheData) {
	val, _ := json.Marshal(*userdata)
	PostRedis(uc, session, val, time.Minute*2+time.Second+2)
}

func (uc *UserCache) PostSessionFlag(session string, userID string) {
	PostRedis(uc, "forget:"+session, userID, time.Minute*5)
}

func (uc *UserCache) GetSessionFlag(session string) (string, error) {
	return GetRedis(uc, "forget:"+session)
}

func (uc *UserCache) PostPhone(phone string, session string) {
	PostRedis(uc, "phone:"+phone, session, time.Minute*2+time.Second*2)
}

func (uc *UserCache) DeleteSession(session string) {
	DeleteRedis(uc, session)
}
