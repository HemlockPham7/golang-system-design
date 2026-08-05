package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/internal/repository/bookmark"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
)

type Service interface {
	CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)
	UpdateBookmark(ctx context.Context, description, url, userID, ID string) (*model.Bookmark, error)
	DeleteBookmark(ctx context.Context, userID, ID string) error
	GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarksResult, error)
}

type bookmarkService struct {
	repo    bookmark.Repository
	codeGen utils.GenPass
}

func NewService(repo bookmark.Repository, codeGen utils.GenPass) Service {
	return &bookmarkService{repo: repo, codeGen: codeGen}
}
