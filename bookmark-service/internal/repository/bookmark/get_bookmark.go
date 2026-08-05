package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/dbutils"
)

func (r *bookmarkRepository) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error) {
	bookmarks, err := r.getBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := r.countBookmarks(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

func (r *bookmarkRepository) getBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error) {
	bookmarks := make([]*model.Bookmark, 0, limit)

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Offset(offset).Limit(limit).Find(&bookmarks).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return bookmarks, nil
}

func (r *bookmarkRepository) countBookmarks(ctx context.Context, userID string) (int64, error) {
	var total int64

	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return 0, dbutils.CatchDBError(err)
	}

	return total, nil
}
