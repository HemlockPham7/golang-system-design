package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/dbutils"
	"gorm.io/gorm/clause"
)

func (r *bookmarkRepository) UpdateBookmarkByID(ctx context.Context, updatedBookmark *model.Bookmark, userID, bookmarkID string) (*model.Bookmark, error) {
	returnedBookmark := &model.Bookmark{}

	result := r.db.WithContext(ctx).
		Model(returnedBookmark).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", bookmarkID, userID).
		Updates(updatedBookmark)

	if result.Error != nil {
		return nil, dbutils.CatchDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, dbutils.ErrRecordNotFoundType
	}

	return returnedBookmark, nil
}
