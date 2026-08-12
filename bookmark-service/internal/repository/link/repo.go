package link

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name Repository --filename repo.go --outpkg mockLink
type Repository interface {
	StoreURL(ctx context.Context, code, url string, expSecond int64) error
	GetURL(ctx context.Context, code string) (string, error)
}

type linkRepository struct {
	c *redis.Client
}

func NewLinkRepository(c *redis.Client) Repository {
	return &linkRepository{c: c}
}
