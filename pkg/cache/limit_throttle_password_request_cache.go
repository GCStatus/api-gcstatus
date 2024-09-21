package cache

import (
	"time"
)

func (r *RedisCache) GetPasswordThrottleCache(key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()

	return value, err
}

func (r *RedisCache) SetPasswordThrottleCache(key string, duration time.Duration) error {
	_, err := r.client.Set(ctx, key, 1, duration).Result()

	return err
}

func (r *RedisCache) RemovePasswordThrottleCache(email string) error {
	key := "password-reset:" + email

	_, err := r.client.Del(ctx, key).Result()

	return err
}
