package service

import (
	"context"
	"testing"

	mockRepo "github.com/HemlockPham7/golang-system-design/internal/repository/mocks"
	mockService "github.com/HemlockPham7/golang-system-design/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockCodeGen func() *mockService.GenPass
		setupMockStorage func(ctx context.Context) *mockRepo.URLStorage

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

			setupMockStorage: func(ctx context.Context) *mockRepo.URLStorage {
				storageMock := mockRepo.NewURLStorage(t)
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

			setupMockStorage: func(ctx context.Context) *mockRepo.URLStorage {
				return mockRepo.NewURLStorage(t)
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

			setupMockStorage: func(ctx context.Context) *mockRepo.URLStorage {
				storageMock := mockRepo.NewURLStorage(t)
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
			testService := NewShortenUrl(mockStorage, mockCodeGen)

			code, err := testService.CreateShortenLink(ctx, tc.inputURL, tc.inputExp)
			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedError, err)

		})
	}
}

func TestGetLinkFromCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockStorage func(ctx context.Context) *mockRepo.URLStorage

		inputCode string

		expectedURL   string
		expectedError error
	}{
		{
			name: "get url successfully",

			setupMockStorage: func(ctx context.Context) *mockRepo.URLStorage {
				storageMock := mockRepo.NewURLStorage(t)
				storageMock.On("GetURL", ctx, "abc1234").Return("google.com", nil)
				return storageMock
			},

			inputCode: "abc1234",

			expectedURL:   "google.com",
			expectedError: nil,
		},
		{
			name: "fail to get code",

			setupMockStorage: func(ctx context.Context) *mockRepo.URLStorage {
				storageMock := mockRepo.NewURLStorage(t)
				storageMock.On("GetURL", ctx, "abc1234").Return("", assert.AnError)
				return storageMock
			},

			inputCode: "abc1234",

			expectedURL:   "",
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockStorage := tc.setupMockStorage(ctx)
			testService := NewShortenUrl(mockStorage, nil)

			url, err := testService.GetLinkFromCode(ctx, tc.inputCode)
			assert.Equal(t, tc.expectedURL, url)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
