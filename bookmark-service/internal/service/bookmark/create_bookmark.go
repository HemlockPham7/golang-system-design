package bookmark

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
)

const codeLength = 8

func (s *bookmarkService) CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	// create code
	code, err := s.codeGen.GeneratePassword(codeLength)
	if err != nil {
		return nil, err
	}

	// create bookmark model
	bookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		Code:        code,
		UserID:      userID,
	}

	// call repo
	res, err := s.repo.CreateBookmark(ctx, bookmark)
	if err != nil {
		return nil, err
	}

	//return
	return res, nil
}
