package users

import (
	"time"
)

type UserProfile struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint
	Bio       string `json:"bio"`
}
