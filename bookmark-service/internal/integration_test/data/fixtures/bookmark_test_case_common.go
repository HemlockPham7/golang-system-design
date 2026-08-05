package fixtures

import (
	"github.com/HemlockPham7/golang-system-design/internal/model"
	"gorm.io/gorm"
)

type BookmarkCommonTestDB struct {
	UserCommonTestDB
}

func (b *BookmarkCommonTestDB) Migrate() error {
	return b.db.AutoMigrate(&model.User{}, &model.Bookmark{})
}

func (b *BookmarkCommonTestDB) GenerateData() error {
	err := b.UserCommonTestDB.GenerateData()
	if err != nil {
		return err
	}

	bookmarks := []*model.Bookmark{
		{
			Base:        GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
			Description: "Bookmark 1",
			URL:         "https://www.google.com",
			Code:        "bookmark1",
			UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
		},
		{
			Base:        GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd7"),
			Description: "Bookmark 2",
			URL:         "https://www.google.com",
			Code:        "bookmark2",
			UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd5",
		},
	}

	return b.db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(bookmarks, 10).Error
}
