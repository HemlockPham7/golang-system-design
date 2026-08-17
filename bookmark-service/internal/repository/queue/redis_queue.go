package queue

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type redisQueue struct {
	client    *redis.Client
	queueName string
}

func NewRedisQueue(c *redis.Client, queueName string) Repository {
	return &redisQueue{
		client:    c,
		queueName: queueName,
	}
}

func (r *redisQueue) PushMessage(ctx context.Context, message []byte) error {
	return r.client.LPush(ctx, r.queueName, message).Err()
}

var NoMessageError = errors.New("no message")

func (r *redisQueue) PopMessage(ctx context.Context) ([]byte, error) {
	msg, err := r.client.RPop(ctx, r.queueName).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, NoMessageError
		}
		return nil, err
	}
	return msg, nil
}
