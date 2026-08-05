package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HemlockPham7/golang-system-design/pkg/utils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("test error")

func TestGenPassHandler_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.GenPass

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)
			},

			setupMockService: func(ctx context.Context) *mocks.GenPass {
				serviceMock := mocks.NewGenPass(t)
				serviceMock.On("GeneratePassword", passwordLength).Return("123456789012", nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"password":"123456789012"}`,
		},
		{
			name: "service failed",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)
			},

			setupMockService: func(ctx context.Context) *mocks.GenPass {
				serviceMock := mocks.NewGenPass(t)
				serviceMock.On("GeneratePassword", passwordLength).Return("", testErr)
				return serviceMock
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Processing Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()        // rec will store http status, response body, and header
			ctx, _ := gin.CreateTestContext(rec) // ctx.Writer = rec, rec is ResponseWriter, means every handler call c.JSON(...) would write the value into rec
			tc.setupRequest(ctx)                 // ctx.Request include request from GET genpass

			mockService := tc.setupMockService(ctx)
			testHandler := NewGenPass(mockService)

			testHandler.GeneratePassword(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
