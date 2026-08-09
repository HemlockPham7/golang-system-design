package queue

import (
	"context"
	"testing"

	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisQueue_PushMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context, t *testing.T) *redis.Client

		inputMessage   []byte
		expectedError  error
		expectedResult []byte
	}{
		{
			name: "success",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				redisClientMock := redisPkg.InitMockRedis(t)
				return redisClientMock
			},
			inputMessage:   []byte("test message"),
			expectedError:  nil,
			expectedResult: []byte("test message"),
		},
		{
			name: "failed due to redis close",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				redisClientMock := redisPkg.InitMockRedis(t)
				_ = redisClientMock.Close()
				return redisClientMock
			},
			inputMessage:  []byte("test message"),
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisClientMock := tc.setupMock(ctx, t)
			storage := NewRedisQueue(redisClientMock, "queue_name")

			err := storage.PushMessage(ctx, tc.inputMessage)
			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, tc.expectedResult, tc.inputMessage)
			}
		})
	}
}
