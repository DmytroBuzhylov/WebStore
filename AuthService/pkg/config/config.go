package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PORT      string
	DB        string
	Redis     string
	JWTSecret string
	RabbitMQ  string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Load config: %v", err)
	}

	AppConfig = &Config{
		PORT:      os.Getenv("PORT"),
		DB:        os.Getenv("DB"),
		Redis:     os.Getenv("REDIS"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		RabbitMQ:  os.Getenv("RABBIT_MQ"),
	}
}
