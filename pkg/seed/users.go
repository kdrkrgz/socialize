package seed

import (
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
