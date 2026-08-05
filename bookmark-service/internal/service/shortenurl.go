package service

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/repository"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
)

const codeLength = 7

//go:generate mockery --name ShortenUrl --filename shortenurl.go --outpkg mocks
type ShortenUrl interface {
	CreateShortenLink(ctx context.Context, url string, expSecond int64) (string, error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

type shortenUrl struct {
	storage repository.URLStorage
	codeGen utils.GenPass
}

func NewShortenUrl(storage repository.URLStorage, codeGen utils.GenPass) ShortenUrl {
	return &shortenUrl{storage: storage, codeGen: codeGen}
}

func (s *shortenUrl) CreateShortenLink(ctx context.Context, url string, expSecond int64) (string, error) {
	// tao code
	code, err := s.codeGen.GeneratePassword(codeLength)
	if err != nil {
		return "", err
	}
	// goi repo de store url
	err = s.storage.StoreURL(ctx, code, url, expSecond)
	if err != nil {
		return "", err
	}
	// return code
	return code, nil
}

func (s *shortenUrl) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return s.storage.GetURL(ctx, code)
}
