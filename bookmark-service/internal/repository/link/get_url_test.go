package link

import (
	"context"
	"testing"

	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestLinkRepository_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context, t *testing.T, code, url string) Repository

		expectedValue string
		expectedError error
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context, t *testing.T, code, url string) Repository {
				mock := redisPkg.InitMockRedis(t)
				storage := NewLinkRepository(mock)
				err := storage.StoreURL(ctx, code, url, 300)
				assert.NoError(t, err)
				return storage
			},

			expectedValue: "google.com",
			expectedError: nil,
		},
		{
			name: "code not found",

			setupMock: func(ctx context.Context, t *testing.T, code, url string) Repository {
				mock := redisPkg.InitMockRedis(t)
				storage := NewLinkRepository(mock)
				return storage
			},

			expectedValue: "",
			expectedError: ErrCodeNotFound,
		},
		{
			name: "redis closed",

			setupMock: func(ctx context.Context, t *testing.T, code, url string) Repository {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				storage := NewLinkRepository(mock)
				return storage
			},

			expectedValue: "",
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			code := "12345"
			url := "google.com"
			storage := tc.setupMock(ctx, t, code, url)

			value, err := storage.GetURL(ctx, code)
			assert.ErrorIs(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedValue, value)
		})
	}
}
