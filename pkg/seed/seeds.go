package seed

import (
	"github.com/kdrkrgz/socalize/utils"
	"gorm.io/gorm"
)

func All() map[string][]Seed {
	hashedPass, _ := utils.HashPassword("pass123")
	seedData := map[string][]Seed{
		"Users": {
			Seed{
				Name: "user1",
				Run: func(db *gorm.DB) error {
					return CreateUser(db, "Jane", "Doe", "jane@doe.com", hashedPass)
				},
			},
			Seed{
				Name: "user2",
				Run: func(db *gorm.DB) error {
					return CreateUser(db, "John", "Doe", "john@doe.com", hashedPass)
				},
			},
			Seed{
				Name: "user1-profile",
				Run: func(db *gorm.DB) error {
					return CreateUserProfile(db, 1, "jane bio")
				},
			},
			Seed{
				Name: "user2-profile",
				Run: func(db *gorm.DB) error {
					return CreateUserProfile(db, 2, "john bio")
				},
			},
			Seed{
				Name: "user1-post1",
				Run: func(db *gorm.DB) error {
					return CreatePost(db, 1, "jane first post")
				},
			},
			Seed{
				Name: "user1-post2",
				Run: func(db *gorm.DB) error {
					return CreatePost(db, 1, "jane second post")
				},
			},
			Seed{
				Name: "user2-post1",
				Run: func(db *gorm.DB) error {
					return CreatePost(db, 2, "john first post")
				},
			},
		},
	}
	return seedData
}
