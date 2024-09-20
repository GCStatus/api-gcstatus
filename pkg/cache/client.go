package cache

import (
	"context"
	"gcstatus/config"

	"github.com/go-redis/redis/v8"
)

var env = config.LoadConfig()
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: env.RedisHost,
})
