package cache

import (
	"context"
	"time"
)

func AddThrottleCache(ctx context.Context, key string) (count int64, error error) {
	count, err := rdb.Incr(ctx, key).Result()

	return count, err
}

func ExpireThrottleCache(ctx context.Context, key string, timeWindow time.Duration) {
	rdb.Expire(ctx, key, timeWindow)
}
