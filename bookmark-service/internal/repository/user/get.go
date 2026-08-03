package user

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/pkg/dbutils"
)

func (r *sqlRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return user, nil
}

func (r *sqlRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return user, nil
}
