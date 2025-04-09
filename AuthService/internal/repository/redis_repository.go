package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisRepository) Del(ctx context.Context, keys ...string) error {
	return r.rdb.Del(ctx, keys...).Err()
}
