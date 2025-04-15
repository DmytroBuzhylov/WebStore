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
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	sqlDB, err := database.InitDB(cfg.DB)
	defer sqlDB.Close()
	db, err := database.InitGorm(sqlDB)

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

	gin.SetMode(gin.DebugMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:              ":" + cfg.PORT,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
