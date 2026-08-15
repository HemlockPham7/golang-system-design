package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/internal/service/queue"
)

func (s *bookmarkService) CreateBatchBookmarks(ctx context.Context, userId string, bookmarkList []*queue.ImportBookmarkInput) error {
	bookmarks := make([]*model.Bookmark, len(bookmarkList))

	for i, input := range bookmarkList {
		code, err := s.codeGen.GeneratePassword(codeLength)
		if err != nil {
			return err
		}

		bookmarks[i] = &model.Bookmark{
			Description: input.Description,
			URL:         input.URL,
			UserID:      userId,
			Code:        code,
		}
	}

	err := s.repo.CreateBatchBookmarks(ctx, bookmarks)
	if err != nil {
		return err
	}
	return nil
}
