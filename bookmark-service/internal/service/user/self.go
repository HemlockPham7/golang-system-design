package user

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
)

func (s *service) GetSelfInfo(ctx context.Context, uid string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, uid)
}
