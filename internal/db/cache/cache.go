package cache

import (
	"context"
	"errors"
	"github.com/hosseintrz/suggestion_api/internal/config"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	ErrUnknownCache = errors.New("unknown cache type")
)

type Cache interface {
	Set(context.Context, string, interface{}, time.Duration) error
	Get(context.Context, string) (string, error)
	Expire(context.Context, string, time.Duration) error
}

func NewCache(conf config.CacheConfig) (Cache, error) {
	switch conf.Type {
	case "redis":
		return NewRedis(conf), nil
	default:
		logrus.Fatal("unknown cache config")
		return nil, ErrUnknownCache
	}
}
