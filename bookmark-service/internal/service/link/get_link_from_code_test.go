package link

import (
	"context"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	mock_bookmark "github.com/HemlockPham7/golang-system-design/internal/repository/bookmark/mocks"
	mockLink "github.com/HemlockPham7/golang-system-design/internal/repository/link/mocks"

	"github.com/stretchr/testify/assert"
)

func TestLinkService_GetLinkFromCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockStorage    func(ctx context.Context) *mockLink.Repository
		setupMockRepository func(ctx context.Context) *mock_bookmark.Repository

		inputCode string

		expectedURL   string
		expectedError error
	}{
		{
			name: "get url successfully",

			setupMockStorage: func(ctx context.Context) *mockLink.Repository {
				storageMock := mockLink.NewRepository(t)
				storageMock.On("GetURL", ctx, "abc1234").Return("google.com", nil)
				return storageMock
			},

			setupMockRepository: func(ctx context.Context) *mock_bookmark.Repository {
				return mock_bookmark.NewRepository(t)
			},

			inputCode: "abc1234",

			expectedURL:   "google.com",
			expectedError: nil,
		},
		{
			name: "fail to get code from GetURL",

			setupMockStorage: func(ctx context.Context) *mockLink.Repository {
				storageMock := mockLink.NewRepository(t)
				storageMock.On("GetURL", ctx, "abc1234").Return("", assert.AnError)
				return storageMock
			},

			setupMockRepository: func(ctx context.Context) *mock_bookmark.Repository {
				return mock_bookmark.NewRepository(t)
			},

			inputCode: "abc1234",

			expectedURL:   "",
			expectedError: assert.AnError,
		},
		{
			name: "successfully to get bookmark from code",

			setupMockStorage: func(ctx context.Context) *mockLink.Repository {
				return mockLink.NewRepository(t)
			},

			setupMockRepository: func(ctx context.Context) *mock_bookmark.Repository {
				mockBookmarkRepository := mock_bookmark.NewRepository(t)
				mockBookmarkRepository.On("GetBookmarkByCode", ctx, "abcd1234").
					Return(&model.Bookmark{
						URL:  "https://www.google.com",
						Code: "abcd1234",
					}, nil)
				return mockBookmarkRepository
			},

			inputCode: "abcd1234",

			expectedURL:   "https://www.google.com",
			expectedError: nil,
		},
		{
			name: "failed to get bookmark from code",

			setupMockStorage: func(ctx context.Context) *mockLink.Repository {
				return mockLink.NewRepository(t)
			},

			setupMockRepository: func(ctx context.Context) *mock_bookmark.Repository {
				mockBookmarkRepository := mock_bookmark.NewRepository(t)
				mockBookmarkRepository.On("GetBookmarkByCode", ctx, "abcd1234").
					Return(nil, assert.AnError)
				return mockBookmarkRepository
			},

			inputCode: "abcd1234",

			expectedURL:   "",
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockStorage := tc.setupMockStorage(ctx)
			mockRepository := tc.setupMockRepository(ctx)
			testService := NewLinkService(mockStorage, mockRepository, nil)

			url, err := testService.GetLinkFromCode(ctx, tc.inputCode)
			assert.Equal(t, tc.expectedURL, url)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
