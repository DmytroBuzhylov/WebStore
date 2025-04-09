package handler

import (
	"AuthService/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/login", h.login)
		api.POST("/register", h.register)
		api.POST("/verify", h.verify)
	}
}
