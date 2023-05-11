package users

import (
	"github.com/go-playground/validator/v10"
	"github.com/kdrkrgz/socalize/posts"
	"time"
)

type User struct {
	Id        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Username  string       `gorm:"unique" json:"username"`
	Email     string       `gorm:"unique" json:"email"`
	Phone     *string      `gorm:"unique;default:null" json:"phone"`
	Password  string       `json:"-"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Role      string       `gorm:"default:user" json:"role"`
	Posts     []posts.Post `gorm:"foreignKey:UserID" json:"posts"`
	Profile   UserProfile  `gorm:"foreignKey:UserID" json:"profile"`
}

type SignUpInput struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	Id        uint      `json:"id,omitempty"`
	FirstName string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FilterUserRecord(u *User) *UserResponse {
	return &UserResponse{
		Id:        u.Id,
		FirstName: u.FirstName,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Value string `json:"value,omitempty"`
	Tag   string `json:"tag"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}
	var errors []*ErrorResponse
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, &ErrorResponse{
			Field: err.StructNamespace(),
			Value: err.Value().(string),
			Tag:   err.Tag(),
		})
	}
	return errors
}
