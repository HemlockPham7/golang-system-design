package link

import (
	"context"
	"errors"

	"github.com/HemlockPham7/golang-system-design/internal/repository/bookmark"
	"github.com/HemlockPham7/golang-system-design/internal/repository/link"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
)

const (
	codeLength         = 7
	codeLengthBookmark = 8
)

var (
	ErrCodeNotFound = errors.New("code not found")
)

//go:generate mockery --name Service --filename service.go --outpkg mockLink
type Service interface {
	CreateShortenLink(ctx context.Context, url string, expSecond int64) (string, error)
	GetLinkFromCode(ctx context.Context, requestCode string) (string, error)
}

type linkService struct {
	linkRepository     link.Repository
	bookmarkRepository bookmark.Repository
	codeGen            utils.GenPass
}

func NewLinkService(linkRepository link.Repository, bookmarkRepository bookmark.Repository, codeGen utils.GenPass) Service {
	return &linkService{linkRepository: linkRepository, bookmarkRepository: bookmarkRepository, codeGen: codeGen}
}
