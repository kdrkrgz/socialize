package users

import (
	p "github.com/kdrkrgz/socalize/posts"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string      `gorm:"unique" json:"username"`
	Email     string      `gorm:"unique" json:"email"`
	Phone     string      `gorm:"unique;default:null" json:"phone"`
	Password  string      `json:"-"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Posts     []p.Post    `gorm:"foreignKey:UserID" json:"posts"`
	Profile   UserProfile `gorm:"foreignKey:UserID" json:"profile"`
}
