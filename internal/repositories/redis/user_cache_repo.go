package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/niflheimdevs/backend/internal/models"
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

func (uc *UserCache) FindByUsername(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	value, err := uc.DB.Get(ctx, "username:"+username).Result()
	return value, err
}

func (uc *UserCache) FindByPhone(phonenumber string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	value, err := uc.DB.Get(ctx, "phone:"+phonenumber).Result()
	return value, err
}

// ? idk if i should break this down
func (uc *UserCache) PostUserCreds(phonenumber string, username string, password []byte, session string, otp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := uc.DB.Set(ctx, "username:"+username, session, time.Minute*2+time.Second*2).Err()
	if err != nil {
		return err
	}

	err = uc.DB.Set(ctx, "phone:"+phonenumber, session, time.Minute*2+time.Second*2).Err()
	if err != nil {
		return err
	}

	userdata := models.UserCacheData{
		Username: username,
		Phone:    phonenumber,
		Password: password,
		OTP:      otp,
	}
	val, _ := json.Marshal(userdata)
	err = uc.DB.Set(ctx, session, val, time.Minute*2+time.Second*2).Err()
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserCache) GetUserCreds(session string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	val, err := uc.DB.Get(ctx, session).Result()
	if err == redis.Nil {
		return ""
	}
	return val
}

func (uc *UserCache) ClearUserCreds(phonenumber string, username string, session string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uc.DB.Del(ctx, "username:"+username)

	uc.DB.Del(ctx, "phone:"+phonenumber)

	uc.DB.Del(ctx, session)
}
func (uc *UserCache) RedisPing() {
	log.Println(uc.DB.Ping(context.Background()).Result())
}
