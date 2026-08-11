package bookmark

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestBookmarkEndpoint_CreateBookmark(t *testing.T) {
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
			name: "create bookmark successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequest("https://google.com", "Google", true)
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				mockJWTValidator := mocks.NewJWTValidator(t)
				mockJWTValidator.On("ValidateJWT", "valid_jwt_token").
					Return(jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"}, nil)
				return mockJWTValidator
			},

			setupRateLimit: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return setupRateLimit(ctx, 1, redisClient)
			},

			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `"message":"Create a bookmark successfully!"`,
		},
		{
			name: "invalid create bookmark payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequest("", "Google", true)
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				mockJWTValidator := mocks.NewJWTValidator(t)
				mockJWTValidator.On("ValidateJWT", "valid_jwt_token").
					Return(jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"}, nil)
				return mockJWTValidator
			},

			setupRateLimit: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return setupRateLimit(ctx, 1, redisClient)
			},

			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `"message":"Input error"`,
		},
		{
			name: "invalid token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequest("https://google.com", "Google", true)
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				mockJWTValidator := mocks.NewJWTValidator(t)
				mockJWTValidator.On("ValidateJWT", "valid_jwt_token").
					Return(nil, assert.AnError)
				return mockJWTValidator
			},

			setupRateLimit: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return setupRateLimit(ctx, 1, redisClient)
			},

			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `"error":"Invalid token"`,
		},
		{
			name: "too many request",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequest("https://google.com", "Google", true)
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				mockJWTValidator := mocks.NewJWTValidator(t)
				mockJWTValidator.On("ValidateJWT", "valid_jwt_token").
					Return(jwt.MapClaims{"sub": "de305d54-75b4-431b-adb2-eb6b9e546099"}, nil)
				return mockJWTValidator
			},

			setupRateLimit: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return setupRateLimit(ctx, middleware.RateLimitCount, redisClient)
			},

			expectedStatusCode:   http.StatusTooManyRequests,
			expectedResponseBody: `"message":"Create a bookmark successfully!"`,
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

func setupRequest(url, description string, haveClaims bool) (*http.Request, *httptest.ResponseRecorder) {
	reqBody := fmt.Sprintf(`{"url":"%s","description":"%s"}`, url, description)
	req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_jwt_token")

	if haveClaims {
		req.Header.Set("claims", "de305d54-75b4-431b-adb2-eb6b9e546099")
	}
	rec := httptest.NewRecorder()
	return req, rec
}

func setupRateLimit(ctx context.Context, rateLimit int, redisClient *redis.Client) *redis.Client {
	key := fmt.Sprintf(middleware.RateLimitKeyFormat, "de305d54-75b4-431b-adb2-eb6b9e546099")
	redisClient.Set(ctx, key, rateLimit, middleware.RateLimitInterval)
	return redisClient
}
