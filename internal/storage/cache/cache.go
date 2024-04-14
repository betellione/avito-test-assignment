package banner

import (
	m "banner/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var Ctx = context.Background()

func FetchBannerFromCache(redisClient *redis.Client, tagID, featureID int) (*m.ResponseBanner, error) {
	key := fmt.Sprintf("banner:%d:%d", tagID, featureID)
	result, err := redisClient.Get(Ctx, key).Result()
	if err != nil {
		return nil, err
	}
	banner := &m.ResponseBanner{}
	if err := json.Unmarshal([]byte(result), banner); err != nil {
		return nil, err
	}
	return banner, nil
}

func CacheBanner(redisClient *redis.Client, tagID, featureID int, banner *m.ResponseBanner) {
	key := fmt.Sprintf("banner:%d:%d", tagID, featureID)
	data, err := json.Marshal(banner)
	if err != nil {
		log.Printf("Error marshalling banner: %v", err)
		return
	}
	err = redisClient.Set(Ctx, key, data, 5*time.Minute).Err()
	if err != nil {
		log.Printf("Error caching banner: %v", err)
	}
}
