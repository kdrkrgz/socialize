package seed

import (
	"errors"
	"github.com/kdrkrgz/socalize/posts"
	"github.com/kdrkrgz/socalize/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitialDataSeed() {
	db := openConnection()
	seeds := All()
	if err := db.AutoMigrate(&users.User{}); err == nil && db.Migrator().HasTable(&users.User{}) {
		if err := db.First(&users.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			runSeeder(db, seeds["InitialData"])
		}
	}

}

func runSeeder(db *gorm.DB, seeds []Seed) {

	for _, seed := range seeds {
		if err := seed.Run(db); err != nil {
			log.Fatalf("Error seeding database: %s err: %s", seed.Name, err)
		}

	}
}

func openConnection() *gorm.DB {
	// Seed the database with some initial data with initial migration
	dsn := "host=0.0.0.0 user=postgres dbname=socialize port=5432 sslmode=disable"

	//dsn := fmt.Sprintf("host:%s user=%s dbname=%s port=%s sslmode=%s",
	//	conf.Get("DbHost"),
	//	conf.Get("DbUser"),
	//	conf.Get("DbName"),
	//	conf.Get("DbPort"),
	//	conf.Get("SSLMode"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(&users.User{}, &posts.Post{}, &users.UserProfile{}); err != nil {
		panic("db migration failed")
	}
	return db
}
