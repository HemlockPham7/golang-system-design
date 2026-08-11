package bookmark

import (
	"context"
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	mock_bookmark "github.com/HemlockPham7/golang-system-design/internal/repository/bookmark/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkService_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context, inputUserID string, inputPage, inputLimit int) *mock_bookmark.Repository

		inputUserID string
		inputPage   int
		inputLimit  int

		expectedError     error
		expectedBookmarks []*model.Bookmark
		expectedTotal     int64
	}{
		{
			name: "success",

			setupMockRepo: func(ctx context.Context, inputUserID string, inputPage, inputLimit int) *mock_bookmark.Repository {
				repoMock := mock_bookmark.NewRepository(t)
				offset := (inputPage - 1) * inputLimit
				repoMock.On("GetBookmarks", ctx, inputUserID, inputLimit, offset).Return(
					[]*model.Bookmark{
						{
							Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
							Description: "Bookmark 1",
							URL:         "https://www.google.com",
							Code:        "bookmark1",
							UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
						},
					},
					int64(1),
					nil,
				)
				return repoMock
			},

			inputUserID: "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			inputPage:   1,
			inputLimit:  10,

			expectedError: nil,
			expectedBookmarks: []*model.Bookmark{
				{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
					Description: "Bookmark 1",
					URL:         "https://www.google.com",
					Code:        "bookmark1",
					UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				},
			},
			expectedTotal: 1,
		},
		{
			name: "get bookmarks failed - repository error",

			setupMockRepo: func(ctx context.Context, inputUserID string, inputPage, inputLimit int) *mock_bookmark.Repository {
				repoMock := mock_bookmark.NewRepository(t)
				offset := (inputPage - 1) * inputLimit
				repoMock.On("GetBookmarks", ctx, inputUserID, inputLimit, offset).Return(nil, int64(0), assert.AnError)
				return repoMock
			},

			inputUserID: "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			inputPage:   1,
			inputLimit:  10,

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repoMock := tc.setupMockRepo(ctx, tc.inputUserID, tc.inputPage, tc.inputLimit)
			service := NewService(repoMock, nil)
			getBookmarksResult, err := service.GetBookmarks(ctx, tc.inputUserID, tc.inputPage, tc.inputLimit)
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedBookmarks, getBookmarksResult.Bookmarks)
				assert.Equal(t, tc.expectedTotal, getBookmarksResult.Total)
			}
		})
	}
}
