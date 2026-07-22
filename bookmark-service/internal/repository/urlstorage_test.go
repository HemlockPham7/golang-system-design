package repository

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
)

func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context, t *testing.T) *redis.Client

		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "12345").Result()
				assert.NoError(t, err)
				assert.Equal(t, "google.com", res)
			},
		},
		{
			name: "connection err",

			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},

			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mock := tc.setupMock(ctx, t)

			storage := NewUrlStorage(mock)
			err := storage.StoreURL(ctx, "12345", "google.com", 300)
			assert.Equal(t, tc.expectedErr, err)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
		})
	}
}

func TestURLStorage_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context, t *testing.T, code, url string) URLStorage

		expectedValue string
		expectedError error
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context, t *testing.T, code, url string) URLStorage {
				mock := redisPkg.InitMockRedis(t)
				storage := NewUrlStorage(mock)
				err := storage.StoreURL(ctx, code, url, 300)
				assert.NoError(t, err)
				return storage
			},

			expectedValue: "google.com",
			expectedError: nil,
		},
		{
			name: "code not found",

			setupMock: func(ctx context.Context, t *testing.T, code, url string) URLStorage {
				mock := redisPkg.InitMockRedis(t)
				storage := NewUrlStorage(mock)
				return storage
			},

			expectedValue: "",
			expectedError: ErrCodeNotFound,
		},
		{
			name: "redis closed",

			setupMock: func(ctx context.Context, t *testing.T, code, url string) URLStorage {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				storage := NewUrlStorage(mock)
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
