package model

type Bookmark struct {
	Base
	Description string `json:"description"`
	URL         string `json:"url"`
	Code        string `json:"code" gorm:"unique"`
	UserID      string `json:"user_id"`
	User        *User  `json:"-" gorm:"references:ID"`
}
