package bookmark

import (
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepository_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(t *testing.T) *gorm.DB

		inputUserID string
		inputPage   int
		inputLimit  int

		expectedError     error
		expectedBookmarks []*model.Bookmark
		expectedTotal     int64
	}{
		{
			name: "success",
			setupMock: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},

			inputUserID: "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			inputPage:   1,
			inputLimit:  10,

			expectedError: nil,
			expectedTotal: 1,
			expectedBookmarks: []*model.Bookmark{
				{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
					Description: "Bookmark 1",
					URL:         "https://www.google.com",
					Code:        "bookmark1",
					UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				},
			},
		},
		{
			name: "record not found",
			setupMock: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},

			inputUserID: "d7c13097-67a7-4eae-a60e-unknown-user",
			inputPage:   1,
			inputLimit:  10,

			expectedError:     nil,
			expectedTotal:     0,
			expectedBookmarks: []*model.Bookmark{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			dbMock := tc.setupMock(t)
			repo := NewRepository(dbMock)

			offset := (tc.inputPage - 1) * tc.inputLimit
			bookmarks, total, err := repo.GetBookmarks(ctx, tc.inputUserID, tc.inputLimit, offset)
			assert.ErrorIs(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedTotal, total)
			assert.Equal(t, tc.expectedBookmarks, bookmarks)
		})
	}
}
