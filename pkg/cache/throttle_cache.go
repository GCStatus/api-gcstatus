package cache

import (
	"time"
)

func (r *RedisCache) AddThrottleCache(key string) (count int64, error error) {
	count, err := r.client.Incr(ctx, key).Result()

	return count, err
}

func (r *RedisCache) ExpireThrottleCache(key string, timeWindow time.Duration) {
	r.client.Expire(ctx, key, timeWindow)
}
