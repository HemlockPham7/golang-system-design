package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Ping interface {
	CheckHealth(ctx context.Context) error
}

type healthCheckRepo struct {
	redistClient *redis.Client
}

func NewPing(redistClient *redis.Client) Ping {
	return &healthCheckRepo{redistClient: redistClient}
}

func (r *healthCheckRepo) CheckHealth(ctx context.Context) error {
	return r.redistClient.Ping(ctx).Err()
}
