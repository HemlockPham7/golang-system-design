package service

import (
	"context"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	mockRepo "github.com/HemlockPham7/golang-system-design/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockPing func(ctx context.Context) *mockRepo.Ping

		expectedResp *model.HealthCheckResponse
		expectedErr  error
	}{
		{
			name: "health check successfully",

			setupMockPing: func(ctx context.Context) *mockRepo.Ping {
				mockPing := mockRepo.NewPing(t)
				mockPing.On("CheckHealth", ctx).Return(nil)
				return mockPing
			},

			expectedResp: &model.HealthCheckResponse{
				Message:     "OK",
				ServiceName: "bookmark-service",
				InstanceID:  "instance-1",
			},
			expectedErr: nil,
		},
		{
			name: "health check failed",

			setupMockPing: func(ctx context.Context) *mockRepo.Ping {
				mockPing := mockRepo.NewPing(t)
				mockPing.
					On("CheckHealth", ctx).
					Return(assert.AnError)

				return mockPing
			},

			expectedResp: nil,
			expectedErr:  assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockPing := tc.setupMockPing(ctx)
			testService := NewHealthCheck("bookmark-service", "instance-1", mockPing)

			resp, err := testService.HealthCheck(ctx)
			assert.Equal(t, tc.expectedResp, resp)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
