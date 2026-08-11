package bookmark

import (
	"context"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	mock_bookmark "github.com/HemlockPham7/golang-system-design/internal/repository/bookmark/mocks"
	mockService "github.com/HemlockPham7/golang-system-design/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkService_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockCodeGen func() *mockService.GenPass
		setupMockStorage func(ctx context.Context) *mock_bookmark.Repository

		inputUserID      string
		inputDescription string
		inputURL         string

		expectedError    error
		expectedBookmark *model.Bookmark
	}{
		{
			name: "create bookmark successfully",

			setupMockCodeGen: func() *mockService.GenPass {
				codeGenMock := mockService.NewGenPass(t)
				codeGenMock.On("GeneratePassword", codeLength).Return("bookmark", nil)
				return codeGenMock
			},

			setupMockStorage: func(ctx context.Context) *mock_bookmark.Repository {
				storageMock := mock_bookmark.NewRepository(t)
				bookmark := &model.Bookmark{
					Description: "Bookmark 1",
					URL:         "https://www.google.com",
					Code:        "bookmark",
					UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				}
				storageMock.On("CreateBookmark", ctx, bookmark).Return(&model.Bookmark{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
					Description: "Bookmark 1",
					URL:         "https://www.google.com",
					Code:        "bookmark",
					UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				},
					nil,
				)
				return storageMock
			},

			inputUserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			inputDescription: "Bookmark 1",
			inputURL:         "https://www.google.com",

			expectedError: nil,
			expectedBookmark: &model.Bookmark{
				Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
				Description: "Bookmark 1",
				URL:         "https://www.google.com",
				Code:        "bookmark",
				UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			},
		},
		{
			name: "fail to create bookmark",

			setupMockCodeGen: func() *mockService.GenPass {
				codeGenMock := mockService.NewGenPass(t)
				codeGenMock.On("GeneratePassword", codeLength).Return("bookmark", nil)
				return codeGenMock
			},

			setupMockStorage: func(ctx context.Context) *mock_bookmark.Repository {
				storageMock := mock_bookmark.NewRepository(t)
				bookmark := &model.Bookmark{
					Description: "Bookmark 1",
					URL:         "https://www.google.com",
					Code:        "bookmark",
					UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				}
				storageMock.On("CreateBookmark", ctx, bookmark).Return(nil, assert.AnError)
				return storageMock
			},

			inputUserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			inputDescription: "Bookmark 1",
			inputURL:         "https://www.google.com",

			expectedError:    assert.AnError,
			expectedBookmark: nil,
		},
		{
			name: "generate password fail",

			setupMockCodeGen: func() *mockService.GenPass {
				codeGenMock := mockService.NewGenPass(t)
				codeGenMock.On("GeneratePassword", codeLength).Return("", assert.AnError)
				return codeGenMock
			},

			setupMockStorage: func(ctx context.Context) *mock_bookmark.Repository {
				return mock_bookmark.NewRepository(t)
			},

			inputUserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			inputDescription: "Bookmark 1",
			inputURL:         "https://www.google.com",

			expectedError:    assert.AnError,
			expectedBookmark: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			codeGenMock := tc.setupMockCodeGen()
			storageMock := tc.setupMockStorage(ctx)
			testService := NewService(storageMock, codeGenMock)

			bookmark, err := testService.CreateBookmark(ctx, tc.inputDescription, tc.inputURL, tc.inputUserID)
			assert.ErrorIs(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedBookmark, bookmark)

		})
	}
}
