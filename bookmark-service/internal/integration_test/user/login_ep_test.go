package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/api"
	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/pkg/jwtutils/mocks"
	redisPkg "github.com/HemlockPham7/golang-system-design/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkEndpoint_GetBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP         func(api api.Engine) *httptest.ResponseRecorder
		setupMockJWTGenerator func(t *testing.T) *mocks.JWTGenerator

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "user login successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequestUserLogin("testuser001", "my_SECURE_password123@")
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				mockJWTGenerator := mocks.NewJWTGenerator(t)
				mockJWTGenerator.On("GenerateJWT", mock.Anything).Return("valid_jwt_token", nil)
				return mockJWTGenerator
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":"valid_jwt_token","message":"Logged in successfully!"}`,
		},
		{
			name: "invalid user login payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequestUserLogin("", "")
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				return mocks.NewJWTGenerator(t)
			},

			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Input error"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setupRedisClient := redisPkg.InitMockRedis(t)
			setupDB := fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			setupJWTGenerator := tc.setupMockJWTGenerator(t)
			testAPI := api.NewEngine(&api.EngineOpts{
				App:         gin.Default(),
				Cfg:         &api.Config{},
				RedisClient: setupRedisClient,
				DbClient:    setupDB,
				JwtGen:      setupJWTGenerator,
			})
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

func setupRequestUserLogin(username, password string) (*http.Request, *httptest.ResponseRecorder) {
	reqBody := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	return req, rec
}
