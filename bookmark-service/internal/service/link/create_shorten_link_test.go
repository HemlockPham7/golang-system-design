package link

import (
	"context"
	"testing"

	mockLink "github.com/HemlockPham7/golang-system-design/internal/repository/link/mocks"
	mockService "github.com/HemlockPham7/golang-system-design/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLinkService_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockCodeGen func() *mockService.GenPass
		setupMockStorage func(ctx context.Context) *mockLink.Repository

		inputURL string
		inputExp int64

		expectedCode  string
		expectedError error
	}{
		{
			name: "shorten url successfully",

			setupMockCodeGen: func() *mockService.GenPass {
				codeGenMock := mockService.NewGenPass(t)
				codeGenMock.On("GeneratePassword", codeLength).Return("abc1234", nil)
				return codeGenMock
			},

			setupMockStorage: func(ctx context.Context) *mockLink.Repository {
				storageMock := mockLink.NewRepository(t)
				storageMock.On("StoreURL", ctx, "abc1234", "google.com", int64(300)).Return(nil)
				return storageMock
			},

			inputURL: "google.com",
			inputExp: 300,

			expectedCode:  "abc1234",
			expectedError: nil,
		},
		{
			name: "Fail to generate code",

			setupMockCodeGen: func() *mockService.GenPass {
				codeGenMock := mockService.NewGenPass(t)
				codeGenMock.On("GeneratePassword", codeLength).Return("", assert.AnError)
				return codeGenMock
			},

			setupMockStorage: func(ctx context.Context) *mockLink.Repository {
				return mockLink.NewRepository(t)
			},

			inputURL: "google.com",
			inputExp: 300,

			expectedCode:  "",
			expectedError: assert.AnError,
		},
		{
			name: "fail to store url",

			setupMockCodeGen: func() *mockService.GenPass {
				codeGenMock := mockService.NewGenPass(t)
				codeGenMock.On("GeneratePassword", codeLength).Return("abc1234", nil)
				return codeGenMock
			},

			setupMockStorage: func(ctx context.Context) *mockLink.Repository {
				storageMock := mockLink.NewRepository(t)
				storageMock.On("StoreURL", ctx, "abc1234", "google.com", int64(300)).Return(assert.AnError)
				return storageMock
			},

			inputURL: "google.com",
			inputExp: 300,

			expectedCode:  "",
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockCodeGen := tc.setupMockCodeGen()
			mockStorage := tc.setupMockStorage(ctx)
			testService := NewLinkService(mockStorage, nil, mockCodeGen)

			code, err := testService.CreateShortenLink(ctx, tc.inputURL, tc.inputExp)
			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedError, err)

		})
	}
}
