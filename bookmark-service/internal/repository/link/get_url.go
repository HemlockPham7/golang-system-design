package link

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrCodeNotFound = errors.New("code not found")

func (s *linkRepository) GetURL(ctx context.Context, code string) (string, error) {
	res, err := s.c.Get(ctx, code).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}
	return res, err
}
