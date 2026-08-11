package link

import (
	"context"
)

func (s *linkService) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	switch {
	case len(code) == codeLength:
		return s.linkRepository.GetURL(ctx, code)
	case len(code) == codeLengthBookmark:
		bookmark, err := s.bookmarkRepository.GetBookmarkByCode(ctx, code)
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	default:
		return "", ErrCodeNotFound
	}

}
