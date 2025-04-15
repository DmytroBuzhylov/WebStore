package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("Error opening database: ", err)
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Println("Error connecting to database: ", err)
		db.Close()
		return nil, err
	}

	log.Println("Database connection pool established")
	return db, nil
}

func InitGorm(db *sql.DB) (*gorm.DB, error) {
	dialect := postgres.New(postgres.Config{
		Conn: db,
	})
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	gormDB, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		log.Printf("Failed to initialize GORM: %v", err)
		return nil, err
	}
	return gormDB, nil
}
