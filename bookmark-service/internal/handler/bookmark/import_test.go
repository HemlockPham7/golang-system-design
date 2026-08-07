package bookmark

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/service/queue"
	"github.com/HemlockPham7/golang-system-design/internal/service/queue/mocks"
	"github.com/HemlockPham7/golang-system-design/pkg/csv"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkHandler_ImportBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest   func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer)
		mockQueueSetup func(ctx context.Context) *mocks.Service
		fileContent    string

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case - success",

			setupRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/test",
					body,
				)
				ctx.Set("claims", jwt.MapClaims{
					"sub": "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				})
				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},

			mockQueueSetup: func(ctx context.Context) *mocks.Service {
				serviceMock := mocks.NewService(t)
				// SendImportBookmarkJob(ctx context.Context, uid string, bookmarkInputs []*ImportBookmarkInput) error
				serviceMock.On(
					"SendImportBookmarkJob",
					mock.Anything,
					"d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
					[]*queue.ImportBookmarkInput{{Description: "Google", URL: "https://www.google.com"}, {Description: "Facebook", URL: "https://www.facebook.com"}}).
					Return(nil)
				return serviceMock
			},

			fileContent:      "description,url\nGoogle,https://www.google.com\nFacebook,https://www.facebook.com",
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"Import bookmark job sent successfully"}`,
		},
		{
			name: "error case - invalid csv",

			setupRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/test",
					body,
				)
				ctx.Set("claims", jwt.MapClaims{
					"sub": "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				})
				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},

			mockQueueSetup: func(ctx context.Context) *mocks.Service {
				serviceMock := mocks.NewService(t)
				return serviceMock
			},

			fileContent:      "description,url\nGoogle,43134134\nFacebook,https://www.facebook.com",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Input error","details":["URL is invalid (url)"]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			writer, body := csv.CreateTestMultipartRequest(t, tc.fileContent)

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx, body, writer)

			handler := NewHandler(nil, tc.mockQueueSetup(ctx))
			handler.ImportBookmarks(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
