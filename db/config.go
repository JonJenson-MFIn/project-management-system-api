package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func Migrate(db *gorm.DB) error {

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&Employee{}); err != nil {
		return err
	}
	return nil
}


func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=aws-0-ap-south-1.pooler.supabase.com user=postgres.lwybefbgqqmvzdzkvqnt password=Joans88@joejon dbname=postgres port=6543 sslmode=disable TimeZone=Asia/Kolkata"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := Migrate(database); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	DB = database
}
