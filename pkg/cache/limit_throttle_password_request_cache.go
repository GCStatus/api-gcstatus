package cache

import (
	"context"
	"time"
)

func GetPasswordThrottleCache(ctx context.Context, key string) (string, error) {
	value, err := rdb.Get(ctx, key).Result()

	return value, err
}

func SetPasswordThrottleCache(ctx context.Context, key string, duration time.Duration) error {
	_, err := rdb.Set(ctx, key, 1, duration).Result()

	return err
}

func RemovePasswordThrottleCache(ctx context.Context, email string) error {
	key := "password-reset:" + email

	_, err := rdb.Del(ctx, key).Result()

	return err
}
