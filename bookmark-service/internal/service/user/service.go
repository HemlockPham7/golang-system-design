package user

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/internal/repository/user"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

type service struct {
	repo   user.Repository
	hasher utils.Hasher
}

func NewService(repo user.Repository, hasher utils.Hasher) Service {
	return &service{repo: repo, hasher: hasher}
}
