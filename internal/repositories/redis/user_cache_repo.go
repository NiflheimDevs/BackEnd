package redis

//requires redis version >= 6 otherwise most package calls will result in syntax error
import (
	"context"
	"encoding/json"
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

// ? extract method
// func (uc *UserCache) PostRedis(key string, value interface{}, duration time.Duration) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	err := uc.DB.Set(ctx, key, value, duration).Err()
// 	return err
// }

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
func (uc *UserCache) PostUserCreds(session string, userdata *models.UserCacheData) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	uc.DB.Set(ctx, "username:"+userdata.Username, session, time.Minute*2+time.Second*2).Err()

	uc.DB.Set(ctx, "phone:"+userdata.Phone, session, time.Minute*2+time.Second*2).Err()

	val, _ := json.Marshal(*userdata)
	uc.DB.Set(ctx, session, val, time.Minute*2+time.Second*2).Err()

}

func (uc *UserCache) FindBySession(session string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	val, err := uc.DB.Get(ctx, session).Result()
	return val, err
}

func (uc *UserCache) ClearUserCreds(phonenumber string, username string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uc.DB.Del(ctx, "username:"+username)

	uc.DB.Del(ctx, "phone:"+phonenumber)

}

func (uc *UserCache) DeleteRow(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uc.DB.Del(ctx, key)
}

func (uc *UserCache) PostSessionOTP(session string, userdata *models.UserCacheData) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	val, _ := json.Marshal(*userdata)
	uc.DB.Set(ctx, session, val, time.Minute*2+time.Second+2).Err()
}

func (uc *UserCache) PostSessionFlag(session string, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uc.DB.Set(ctx, "forget:"+session, userID, time.Minute*5).Err()
}

func (uc *UserCache) GetSessionFlag(session string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	value, err := uc.DB.Get(ctx, "forget:"+session).Result()
	return value, err
}
