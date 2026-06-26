package repository

import (
	"context"
	"sociul-auth-service/internal/sentinel"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redisCache *redis.Client
}

func NewCache(redisCache *redis.Client) *Cache {
	return &Cache{redisCache: redisCache}
}

func (r *Cache) Store(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.redisCache.Set(ctx, key, value, ttl).Err()
}

func (r *Cache) Fetch(ctx context.Context, key string) (string, error) {
	value, err := r.redisCache.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", sentinel.ErrCacheMiss
	} else if err != nil {
		return "", err
	}
	return value, nil
}

func (r *Cache) Delete(ctx context.Context, key string) error {
	return r.redisCache.Del(ctx, key).Err()
}

func (r *Cache) Incr(ctx context.Context, key string) error {
	value, err := r.redisCache.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	// Set a 15-minute ttl
	if value == 1 {
		return r.redisCache.Expire(ctx, key, 15*time.Minute).Err()
	}
	return nil
}
