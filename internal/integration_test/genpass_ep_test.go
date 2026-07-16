package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestGenPassEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "normal case",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/genpass", nil)
				respRecorder := httptest.NewRecorder()

				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"password":`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testAPI := api.NewEngine()
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}
