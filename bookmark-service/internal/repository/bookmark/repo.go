package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"gorm.io/gorm"
)

//go:generate mockery --name Repository --filename repo.go --outpkg mock_bookmark
type Repository interface {
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error)
	UpdateBookmarkByID(ctx context.Context, updatedBookmark *model.Bookmark, userID, bookmarkID string) (*model.Bookmark, error)
	DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error
	GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error)
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &bookmarkRepository{db: db}
}
