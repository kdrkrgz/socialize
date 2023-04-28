package users

import (
	p "github.com/kdrkrgz/socalize/posts"
	"time"
)

type User struct {
	ID        uint        `gorm:"primarykey" json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Username  string      `gorm:"unique" json:"username"`
	Email     string      `gorm:"unique" json:"email"`
	Phone     *string     `gorm:"unique;default:null" json:"phone"`
	Password  string      `json:"-"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Posts     []p.Post    `gorm:"foreignKey:UserID" json:"posts"`
	Profile   UserProfile `gorm:"foreignKey:UserID" json:"profile"`
}
