package repository

import (
	"context"
	"testing"

	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestPing_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupMock     func(ctx context.Context, t *testing.T) *redis.Client
		expectedError error
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedError: nil,
		},
		{
			name: "redis closed",
			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mock := tc.setupMock(ctx, t)
			storage := NewPing(mock)

			err := storage.CheckHealth(ctx)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
