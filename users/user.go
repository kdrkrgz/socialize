package users

import (
	"github.com/google/uuid"
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

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
