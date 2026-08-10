package bookmark

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/api"
	"github.com/HemlockPham7/golang-system-design/internal/api/middleware"
	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/pkg/jwtutils/mocks"
	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkEndpoint_GetBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP         func(api api.Engine) *httptest.ResponseRecorder
		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator
		setupRateLimit        func(ctx context.Context, redisClient *redis.Client) *redis.Client

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "get bookmark successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequestGetBookmark(true)
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				mockJWTValidator := mocks.NewJWTValidator(t)
				mockJWTValidator.On("ValidateJWT", "valid_jwt_token").
					Return(jwt.MapClaims{"sub": "d7c13097-67a7-4eae-a60e-0b9b533b7bd4"}, nil)
				return mockJWTValidator
			},

			setupRateLimit: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return setupRateLimitGetBookmark(ctx, 1, redisClient)
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":[{"id":"d7c13097-67a7-4eae-a60e-0b9b533b7bd6","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","description":"Bookmark 1","url":"https://www.google.com","code":"bookmark1","user_id":"d7c13097-67a7-4eae-a60e-0b9b533b7bd4"}],"pagination":{"page":1,"limit":2,"total":1}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setupRedisClient := redisPkg.InitMockRedis(t)
			redisClient := tc.setupRateLimit(context.Background(), setupRedisClient)
			setupDB := fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			setupJWTValidator := tc.setupMockJWTValidator(t)
			testAPI := api.NewEngine(&api.EngineOpts{
				App:         gin.Default(),
				Cfg:         &api.Config{},
				RedisClient: redisClient,
				DbClient:    setupDB,
				JwtVal:      setupJWTValidator,
			})
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

func setupRequestGetBookmark(haveClaims bool) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks/?page=1&limit=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_jwt_token")

	if haveClaims {
		req.Header.Set("claims", "d7c13097-67a7-4eae-a60e-0b9b533b7bd4")
	}
	rec := httptest.NewRecorder()
	return req, rec
}

func setupRateLimitGetBookmark(ctx context.Context, rateLimit int, redisClient *redis.Client) *redis.Client {
	key := fmt.Sprintf(middleware.RateLimitKeyFormat, "d7c13097-67a7-4eae-a60e-0b9b533b7bd4")
	redisClient.Set(ctx, key, rateLimit, middleware.RateLimitInterval)
	return redisClient
}
