package cache

import (
	"context"
	"gcstatus/config"
	"gcstatus/internal/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	AddThrottleCache(key string) (int64, error)
	ExpireThrottleCache(key string, timeWindow time.Duration)
	GetPasswordThrottleCache(key string) (string, error)
	SetPasswordThrottleCache(key string, duration time.Duration) error
	RemovePasswordThrottleCache(email string) error
	GetUserFromCache(userID uint) (*domain.User, bool)
	SetUserInCache(user *domain.User)
	RemoveUserFromCache(userID uint)
}

type RedisCache struct {
	client *redis.Client
}

var ctx = context.Background()
var env = config.LoadConfig()

func NewRedisCache() *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: env.RedisHost,
	})

	return &RedisCache{client: rdb}
}

var GlobalCache Cache
