package redis

import "github.com/redis/go-redis/v9"

type UserCache struct {
	DB *redis.Client
}

func NewUserCache(DB *redis.Client) *UserCache {
	return &UserCache{
		DB: DB,
	}
}
