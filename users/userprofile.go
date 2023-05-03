package users

import (
	"time"
)

type UserProfile struct {
	Id        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint      `json:"userId"`
	Bio       string    `json:"bio"`
}
