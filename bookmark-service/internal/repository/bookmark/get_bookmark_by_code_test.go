package bookmark

import (
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepository_GetBookmarkByCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(t *testing.T) *gorm.DB

		inputCode string

		expectedError    error
		expectedBookmark *model.Bookmark
	}{
		{
			name: "found one bookmark successfully",

			setupMock: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},

			inputCode:     "bookmark1",
			expectedError: nil,
			expectedBookmark: &model.Bookmark{
				Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
				Description: "Bookmark 1",
				URL:         "https://www.google.com",
				Code:        "bookmark1",
				UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			},
		},
		{
			name: "not found one bookmark with input code",

			setupMock: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},

			inputCode:        "not_found",
			expectedError:    dbutils.ErrRecordNotFound,
			expectedBookmark: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			dbMock := tc.setupMock(t)
			repo := NewRepository(dbMock)

			bookmark, err := repo.GetBookmarkByCode(ctx, tc.inputCode)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedBookmark, bookmark)
		})
	}
}
