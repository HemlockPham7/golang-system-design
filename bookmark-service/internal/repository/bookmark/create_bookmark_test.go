package bookmark

import (
	"testing"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepository_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB
		input   *model.Bookmark

		expectErr error
		verifyFn  func(db *gorm.DB, expectedBookmark *model.Bookmark)
	}{
		{
			name: "success",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			input: &model.Bookmark{
				Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd8"),
				Description: "Bookmark 8",
				URL:         "https://www.google.com",
				Code:        "99999999",
				UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			},

			expectErr: nil,
			verifyFn: func(db *gorm.DB, expectedBookmark *model.Bookmark) {
				var bookmark *model.Bookmark
				err := db.Where("code = ?", "99999999").First(&bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, expectedBookmark.ID, bookmark.ID)
				assert.Equal(t, expectedBookmark.Description, bookmark.Description)
				assert.Equal(t, expectedBookmark.Code, bookmark.Code)
				assert.Equal(t, expectedBookmark.URL, bookmark.URL)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)

			bookmark, err := repo.CreateBookmark(ctx, tc.input)
			assert.ErrorIs(t, tc.expectErr, err)
			tc.verifyFn(db, bookmark)
		})
	}
}
