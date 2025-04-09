package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connect to database Error")
	}
	return db
}
