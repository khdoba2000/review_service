package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"                     // db automigration
	_ "github.com/golang-migrate/migrate/v4/database"          // db automigration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // db automigration
	_ "github.com/golang-migrate/migrate/v4/source/file"       // db automigration
	_ "github.com/joho/godotenv/autoload"                      // load .env file automatically
	"go.uber.org/zap"

	_ "github.com/lib/pq" // db driver

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init initializes database connection, migrates then connect with postgres
func Init(dbURL string) (*gorm.DB, error) {

	fmt.Println("init db")

	m, err := migrate.New("file://src/review_service/migrations", dbURL)
	if err != nil {
		log.Fatal("error in creating migrations: ", zap.Error(err))
	}
	fmt.Printf("")
	if err = m.Up(); err != nil {
		log.Println("error updating migrations: ", zap.Error(err))
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	return db, nil
}
