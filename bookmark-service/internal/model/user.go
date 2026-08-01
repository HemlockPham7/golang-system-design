package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string `json:"id" gorm:"type:uuid;primaryKey"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username" gorm:"unique"`
	Password    string `json:"-"`
	Email       string `json:"email" gorm:"unique"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
