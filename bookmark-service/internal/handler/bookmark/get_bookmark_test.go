package bookmark

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/internal/service/bookmark"
	mock_bookmark "github.com/HemlockPham7/golang-system-design/internal/service/bookmark/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkHandler_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mock_bookmark.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful list bookmarks",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Set("claims", jwt.MapClaims{"sub": "user-123"})
			},

			setupMockSvc: func(ctx *gin.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("GetBookmarks", ctx, "user-123", 1, 2).
					Return(&bookmark.GetBookmarksResult{
						[]*model.Bookmark{
							{
								Base:        fixtures.GetTestBase("bookmark-456"),
								Description: "bookmark description",
								URL:         "bookmark url",
								Code:        "bookmark",
								UserID:      "user-123",
							},
						},
						int64(1),
					}, nil)
				return mockService
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"data":[{"id":"bookmark-456","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","description":"bookmark description","url":"bookmark url","code":"bookmark","user_id":"user-123"}],"pagination":{"page":1,"limit":2,"total":1}}`,
		},
		{
			name: "unauthorized user",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context) *mock_bookmark.Service {
				return mock_bookmark.NewService(t)
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"claim not exist"}`,
		},
		{
			name: "internal service",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest("GET", "/v1/bookmarks?page=1&limit=2", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Set("claims", jwt.MapClaims{"sub": "user-123"})
			},

			setupMockSvc: func(ctx *gin.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("GetBookmarks", ctx, "user-123", 1, 2).
					Return(nil, assert.AnError)
				return mockService
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Processing Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			serviceMock := tc.setupMockSvc(ctx)

			handler := NewHandler(serviceMock, nil)

			handler.GetBookmarks(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
