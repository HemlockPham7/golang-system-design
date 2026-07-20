package service

import (
	"context"
	"time"

	"github.com/HemlockPham7/golang-system-design/internal/repository"
)

const codeLength = 7

//go:generate mockery --name ShortenUrl --filename shortenurl.go --outpkg mocks
type ShortenUrl interface {
	CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

type shortenUrl struct {
	storage repository.URLStorage
	codeGen GenPass
}

func NewShortenUrl(storage repository.URLStorage, codeGen GenPass) ShortenUrl {
	return &shortenUrl{storage: storage, codeGen: codeGen}
}

func (s *shortenUrl) CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error) {
	// tao code
	code, err := s.codeGen.GeneratePassword(codeLength)
	if err != nil {
		return "", err
	}
	// goi repo de store url
	err = s.storage.StoreURL(ctx, code, url, exp)
	if err != nil {
		return "", err
	}
	// return code
	return code, nil
}

func (s *shortenUrl) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return s.storage.GetURL(ctx, code)
}
