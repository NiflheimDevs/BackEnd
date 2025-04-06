package redis

//requires redis version >= 6 otherwise most package calls will result in syntax error
import (
	"context"
	"encoding/json"
	"time"

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

func (uc *UserCache) PostRedis(key string, value interface{}, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	uc.DB.Set(ctx, key, value, duration)
}
func (uc *UserCache) DeleteRedis(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uc.DB.Del(ctx, key)
}
func (uc *UserCache) GetRedis(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return uc.DB.Get(ctx, key).Result()
}

func (uc *UserCache) FindByUsername(username string) (string, error) {

	return uc.GetRedis("username:" + username)
}

func (uc *UserCache) FindByPhone(phonenumber string) (string, error) {

	return uc.GetRedis("phone:" + phonenumber)
}

// ? idk if i should break this down
func (uc *UserCache) PostUserCreds(session string, userdata *models.UserCacheData) {

	uc.PostRedis("username:"+userdata.Username, session, time.Minute*2+time.Second*2)
	uc.PostRedis("phone:"+userdata.Phone, session, time.Minute*2+time.Second*2)

	val, _ := json.Marshal(*userdata)
	uc.PostRedis(session, val, time.Minute*2+time.Second*2)

}

func (uc *UserCache) FindBySession(session string) (string, error) {

	return uc.GetRedis(session)
}

func (uc *UserCache) ClearUserCreds(phonenumber string, username string) {
	uc.DeleteRedis("username:" + username)

	uc.DeleteRedis("phone:" + phonenumber)

}

func (uc *UserCache) PostSessionOTP(session string, userdata *models.UserCacheData) {
	val, _ := json.Marshal(*userdata)
	uc.PostRedis(session, val, time.Minute*2+time.Second+2)
}

func (uc *UserCache) PostSessionFlag(session string, userID string) {
	uc.PostRedis("forget:"+session, userID, time.Minute*5)
}

func (uc *UserCache) GetSessionFlag(session string) (string, error) {
	return uc.GetRedis("forget:" + session)
}

func (uc *UserCache) PostPhone(phone string, session string) {
	uc.PostRedis("phone:"+phone, session, time.Minute*2+time.Second*2)
}
