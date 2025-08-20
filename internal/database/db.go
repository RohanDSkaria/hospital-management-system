package database

import (
	"fmt"
	"log"
	"os"

	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection pool
var DB *gorm.DB

// Connect initializes the database connection and runs migrations
func Connect() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	fmt.Println("Staring automigration...")
	err = DB.AutoMigrate(&model.User{}, &model.Patient{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	fmt.Println("Database migration successful!")
}
