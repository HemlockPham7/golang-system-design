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
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &bookmarkRepository{db: db}
}
