package cache

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/hosseintrz/suggestion_api/internal/config"
	"time"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(conf config.CacheConfig) Cache {
	return &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:            conf.Address,
			Username:        conf.Username,
			Password:        conf.PASSWORD,
			DB:              conf.DB,
			ConnMaxIdleTime: conf.ConnMaxIdleTime,
		}),
	}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	return err
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	res, err := r.Client.Get(ctx, key).Result()
	return res, err
}

func (r *Redis) Expire(ctx context.Context, key string, exp time.Duration) error {
	_, err := r.Client.Expire(ctx, key, exp).Result()
	return err
}
