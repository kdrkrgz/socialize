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
				Name: "CreateJane",
				Run: func(db *gorm.DB) error {
					return CreateUser(db, "Jane", "Doe", "jane@doe.com", hashedPass)
				},
			},

			Seed{
				Name: "CreateJohn",
				Run: func(db *gorm.DB) error {
					return CreateUser(db, "John", "Doe", "john@doe.com", hashedPass)
				},
			},
		},
	}
	return seedData
	//return []Seed{
	//	// User seed
	//	Seed{
	//		Name: "CreateJane",
	//		Run: func(db *gorm.DB) error {
	//			return CreateUser(db, "Jane", "Doe", "jane@doe.com", hashedPass)
	//		},
	//	},
	//
	//	Seed{
	//		Name: "CreateJohn",
	//		Run: func(db *gorm.DB) error {
	//			return CreateUser(db, "John", "Doe", "john@doe.com", "pass123")
	//		},
	//	},
	//}
}
