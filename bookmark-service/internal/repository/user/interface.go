package user

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}
