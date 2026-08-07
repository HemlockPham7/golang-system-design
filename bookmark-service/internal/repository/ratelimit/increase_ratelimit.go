package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (r *redisRepo) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) {
	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, exp)
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to increase rate limit")
	}
}
