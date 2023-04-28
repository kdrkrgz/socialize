package posts

import (
	"time"
)

type Post struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint      `gorm:"not null" json:"userId"`
	Content   string    `gorm:"type:text;not null" json:"content"`
}
