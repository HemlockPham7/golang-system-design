package link

import "context"

func (s *linkService) CreateShortenLink(ctx context.Context, url string, expSecond int64) (string, error) {
	// tao code
	code, err := s.codeGen.GeneratePassword(codeLength)
	if err != nil {
		return "", err
	}
	// goi repo de store url
	err = s.linkRepository.StoreURL(ctx, code, url, expSecond)
	if err != nil {
		return "", err
	}
	// return code
	return code, nil
}
