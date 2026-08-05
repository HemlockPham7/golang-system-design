package cache

import (
	"context"
	"time"
)

//go:generate mockery --name DB --filename common.go --outpkg mock_cache
type DB interface {
	SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error
	GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error)
	DeleteCache(ctx context.Context, key string) error
}
