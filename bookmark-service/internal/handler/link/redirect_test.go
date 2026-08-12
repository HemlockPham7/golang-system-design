package link

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/service/link"
	mockLink "github.com/HemlockPham7/golang-system-design/internal/service/link/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLinkHandler_Redirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mockLink.Service

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{"code", "1234567"}}
			},

			setupMockService: func(ctx context.Context) *mockLink.Service {
				serviceMock := mockLink.NewService(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("https://www.google.com", nil)
				return serviceMock
			},

			expectedStatus:   http.StatusFound,
			expectedResponse: "https://www.google.com",
		},
		{
			name: "code is empty",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{"code", ""}}
			},

			setupMockService: func(ctx context.Context) *mockLink.Service {
				return mockLink.NewService(t)
			},

			expectedStatus:   http.StatusBadRequest,
			expectedResponse: ``,
		},
		{
			name: "code not found",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{"code", "1234567"}}
			},

			setupMockService: func(ctx context.Context) *mockLink.Service {
				serviceMock := mockLink.NewService(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("", link.ErrCodeNotFound)
				return serviceMock
			},

			expectedStatus:   http.StatusNotFound,
			expectedResponse: ``,
		},
		{
			name: "internal error server",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{"code", "1234567"}}
			},

			setupMockService: func(ctx context.Context) *mockLink.Service {
				serviceMock := mockLink.NewService(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("", assert.AnError)
				return serviceMock
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockService := tc.setupMockService(ctx)
			testHandler := NewLinkHandler(mockService)

			testHandler.Redirect(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Header().Get("Location"))
		})
	}
}
