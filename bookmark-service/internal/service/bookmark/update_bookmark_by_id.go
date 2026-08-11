package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
)

func (s *bookmarkService) UpdateBookmarkByID(ctx context.Context, description, url, uid, bookmarkID string) (*model.Bookmark, error) {
	updatedBookmark := &model.Bookmark{
		Description: description,
		URL:         url,
	}
	return s.repo.UpdateBookmarkByID(ctx, updatedBookmark, uid, bookmarkID)
}
