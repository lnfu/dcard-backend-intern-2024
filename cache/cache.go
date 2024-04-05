package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	db "github.com/lnfu/dcard-intern/db/sqlc"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	redisClient *redis.Client
}

func NewCache(addr string, password string, db int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis: 無法連接 (%v)\n", err)
	}
	log.Printf("Redis: %v\n", pong)

	return &Cache{client}
}

func generateGetAdvertisementsCacheKey(params db.GetActiveAdvertisementsParams) string {
	components := make([]string, 0)
	if params.Age.Valid {
		components = append(components, fmt.Sprintf("age:%d", params.Age.Int32))
	}
	if params.Gender.Valid {
		components = append(components, fmt.Sprintf("gender:%s", params.Gender.String))
	}
	if params.Country.Valid {
		components = append(components, fmt.Sprintf("country:%s", params.Country.String))
	}
	if params.Platform.Valid {
		components = append(components, fmt.Sprintf("platform:%s", params.Platform.String))
	}
	components = append(components,
		fmt.Sprintf("offset:%d", params.Offset),
		fmt.Sprintf("limit:%d", params.Limit),
	)
	return strings.Join(components, "|")
}

func (cache *Cache) GetAdvertisementsFromCache(ctx context.Context, params db.GetActiveAdvertisementsParams) ([]db.Advertisement, error) {
	key := generateGetAdvertisementsCacheKey(params)
	val, err := cache.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var ads []db.Advertisement
	if err := json.Unmarshal([]byte(val), &ads); err != nil {
		return nil, err
	}
	return ads, nil
}

func (cache *Cache) SetAdvertisementsToCache(ctx context.Context, params db.GetActiveAdvertisementsParams, ads []db.Advertisement) error {
	jsonData, err := json.Marshal(ads)
	if err != nil {
		return err
	}
	key := generateGetAdvertisementsCacheKey(params)
	err = cache.redisClient.Set(ctx, key, jsonData, time.Minute*5).Err()
	if err != nil {
		return err
	}
	return nil
}
