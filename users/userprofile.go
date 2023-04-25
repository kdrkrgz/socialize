package users

import "gorm.io/gorm"

type UserProfile struct {
	gorm.Model
	UserID uint
	Bio    string `json:"bio"`
}
