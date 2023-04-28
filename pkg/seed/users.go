package seed

import (
	"github.com/kdrkrgz/socalize/posts"
	"github.com/kdrkrgz/socalize/users"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, firstName string, lastName string, email string, password string) error {
	return db.Create(&users.User{
		FirstName: firstName,
		LastName:  lastName,
		Username:  firstName + lastName,
		Email:     email,
		Password:  password,
	}).Error
}

func CreateUserProfile(db *gorm.DB, userId uint, bio string) error {
	return db.Create(&users.UserProfile{
		UserID: userId,
		Bio:    bio,
	}).Error
}

func CreatePost(db *gorm.DB, userId uint, content string) error {
	return db.Create(&posts.Post{
		UserID:  userId,
		Content: content,
	}).Error
}
