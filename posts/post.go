package posts

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID  uint   `gorm:"not null" json:"user_id"`
	Content string `gorm:"type:text;not null" json:"content"`
}
