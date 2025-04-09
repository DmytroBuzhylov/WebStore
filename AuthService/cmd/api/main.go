package main

import (
	"AuthService/internal/broker"
	"AuthService/internal/broker/producer"
	"AuthService/internal/database"
	"AuthService/internal/domain"
	"AuthService/internal/handler"
	"AuthService/internal/repository"
	"AuthService/internal/service"
	"AuthService/pkg/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config.LoadConfig()

	cfg := config.AppConfig
	db := database.ConnectPostgres(cfg.DB)
	redis := database.ConnectRedis(cfg.Redis)
	redisRepo := repository.NewRedisRepository(redis)
	rmqConn, err := broker.ConnectRabbit(cfg.RabbitMQ)
	rmq := broker.NewRabbitMQ(rmqConn)
	Producer := producer.NewProducer(rmq)
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, redisRepo, Producer)
	authHandler := handler.NewAuthHandler(authService)

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("error migrrate database: %v", err)
	}

	r := gin.Default()

	authHandler.RegisterRoutes(r)

	r.Run(":" + cfg.PORT)
}
