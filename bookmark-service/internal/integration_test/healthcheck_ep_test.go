package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/api"
	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "health check successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodGet,
					"/health-check",
					nil,
				)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)
				return rec
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"message":"OK"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setupRedisClient := redisPkg.InitMockRedis(t)
			testAPI := api.NewEngine(&api.Config{}, setupRedisClient)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}
