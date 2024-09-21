package cache

import (
	"encoding/json"
	"fmt"
	"gcstatus/internal/domain"
	"log"

	"github.com/go-redis/redis/v8"
)

func (r *RedisCache) GetUserFromCache(userID uint) (*domain.User, bool) {
	key := fmt.Sprintf("user:%d", userID)
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false
		}
		log.Println("Redis error:", err)
		return nil, false
	}

	var user domain.User
	if err := json.Unmarshal([]byte(result), &user); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, false
	}
	return &user, true
}

func (r *RedisCache) SetUserInCache(user *domain.User) {
	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}
	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		log.Println("Redis error:", err)
	}
}

func (r *RedisCache) RemoveUserFromCache(userID uint) {
	key := fmt.Sprintf("user:%d", userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Println("Redis error:", err)
	}
}
