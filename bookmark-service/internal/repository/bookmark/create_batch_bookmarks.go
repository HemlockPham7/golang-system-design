package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/dbutils"
	"gorm.io/gorm"
)

func (r *bookmarkRepository) CreateBatchBookmarks(ctx context.Context, bookmarks []*model.Bookmark) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, bookmark := range bookmarks {
			if err := tx.Create(bookmark).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return dbutils.CatchDBError(err)
	}
	return nil
}
