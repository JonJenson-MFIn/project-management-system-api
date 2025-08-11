package db

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := database.AutoMigrate(
		&Employee{},
		&Project{},
		&Task{},
		&Team{},
		&Ticket{},
		&Notification{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	DB = database
}
