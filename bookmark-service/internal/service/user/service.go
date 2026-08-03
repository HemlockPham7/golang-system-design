package user

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/internal/repository/user"
	"github.com/HemlockPham7/golang-system-design/pkg/jwtutils"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetSelfInfo(ctx context.Context, uid string) (*model.User, error)
}

type service struct {
	repo         user.Repository
	hasher       utils.Hasher
	jwtGenerator jwtutils.JWTGenerator
}

func NewService(repo user.Repository, hasher utils.Hasher, jwtGen jwtutils.JWTGenerator) Service {
	return &service{repo: repo, hasher: hasher, jwtGenerator: jwtGen}
}
