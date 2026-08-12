package link

import (
	"context"
)

func (s *linkService) GetLinkFromCode(ctx context.Context, requestCode string) (string, error) {
	switch {
	case len(requestCode) == codeLength:
		return s.linkRepository.GetURL(ctx, requestCode)
	case len(requestCode) == codeLengthBookmark:
		bookmark, err := s.bookmarkRepository.GetBookmarkByCode(ctx, requestCode)
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	default:
		return "", ErrCodeNotFound
	}
}
