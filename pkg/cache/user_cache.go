package cache

import (
	"encoding/json"
	"fmt"
	"gcstatus/internal/domain"
	"log"

	"github.com/go-redis/redis/v8"
)

func GetUserFromCache(userID uint) (*domain.User, bool) {
	key := fmt.Sprintf("user:%d", userID)
	result, err := rdb.Get(ctx, key).Result()
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

func SetUserInCache(userID uint, user *domain.User) {
	key := fmt.Sprintf("user:%d", userID)
	data, err := json.Marshal(user)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}
	if err := rdb.Set(ctx, key, data, 0).Err(); err != nil {
		log.Println("Redis error:", err)
	}
}

func RemoveUserFromCache(userID uint) {
	key := fmt.Sprintf("user:%d", userID)
	if err := rdb.Del(ctx, key).Err(); err != nil {
		log.Println("Redis error:", err)
	}
}
