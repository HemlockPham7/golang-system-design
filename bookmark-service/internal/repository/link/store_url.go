package link

import (
	"context"
	"time"
)

func (s *linkRepository) StoreURL(ctx context.Context, code, url string, expSecond int64) error {
	return s.c.Set(ctx, code, url, time.Duration(expSecond)*time.Second).Err()
}
