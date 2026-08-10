package bookmark

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	mock_bookmark "github.com/HemlockPham7/golang-system-design/internal/service/bookmark/mocks"
	"github.com/HemlockPham7/golang-system-design/pkg/dbutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var jwtClaims = jwt.MapClaims{
	"sub": "user-123",
}

func TestBookmarkHandler_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *createBookmarkRequest

		setupRequest func(ctx *gin.Context, inputRequest *createBookmarkRequest)
		setupMockSvc func(ctx *gin.Context, inputRequest *createBookmarkRequest) *mock_bookmark.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful create bookmark",

			inputRequest: &createBookmarkRequest{
				URL:         "https://example.com",
				Description: "Example website",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookamrks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Set("claims", jwtClaims)
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *mock_bookmark.Service {
				serviceMock := mock_bookmark.NewService(t)
				serviceMock.On("CreateBookmark", ctx, inputRequest.Description, inputRequest.URL, "user-123").
					Return(&model.Bookmark{
						Base:        fixtures.GetTestBase("bookmark-456"),
						Description: inputRequest.Description,
						URL:         inputRequest.URL,
						Code:        "bookmark",
						UserID:      "user-123",
					}, nil)
				return serviceMock
			},

			expectedCode:     http.StatusCreated,
			expectedResponse: `{"data":{"id":"bookmark-456","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","description":"Example website","url":"https://example.com","code":"bookmark","user_id":"user-123"},"message":"Create a bookmark successfully!"}`,
		},
		{
			name: "unauthorized user",

			inputRequest: &createBookmarkRequest{
				URL:         "https://example.com",
				Description: "Example website",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookamrks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *mock_bookmark.Service {
				return mock_bookmark.NewService(t)
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"claim not exist"}`,
		},
		{
			name: "user id invalid - not exist in system",

			inputRequest: &createBookmarkRequest{
				URL:         "https://example.com",
				Description: "Example website",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookamrks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Set("claims", jwtClaims)
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *mock_bookmark.Service {
				serviceMock := mock_bookmark.NewService(t)
				serviceMock.On("CreateBookmark", ctx, inputRequest.Description, inputRequest.URL, "user-123").
					Return(nil, dbutils.ErrForeignKeyType)
				return serviceMock
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid request",

			inputRequest: &createBookmarkRequest{
				URL:         "",
				Description: "",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookamrks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Set("claims", jwtClaims)
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *mock_bookmark.Service {
				return mock_bookmark.NewService(t)
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Input error","details":["URL is invalid (required)"]}`,
		},
		{
			name: "service internal error",

			inputRequest: &createBookmarkRequest{
				URL:         "https://youtube.com",
				Description: "Webiste to watch videos",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookamrks", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Set("claims", jwtClaims)
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createBookmarkRequest) *mock_bookmark.Service {
				serviceMock := mock_bookmark.NewService(t)
				serviceMock.On("CreateBookmark", ctx, inputRequest.Description, inputRequest.URL, "user-123").
					Return(nil, assert.AnError)
				return serviceMock
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

			tc.setupRequest(ctx, tc.inputRequest)
			serviceMock := tc.setupMockSvc(ctx, tc.inputRequest)

			handler := &bookmarkHandler{
				bookmarkService: serviceMock,
			}

			handler.CreateBookmark(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
