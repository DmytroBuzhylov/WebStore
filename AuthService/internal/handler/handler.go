package handler

import (
	"AuthService/internal/domain"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (h *AuthHandler) register(c *gin.Context) {
	var req *domain.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
		return
	}

	if err := h.authService.Register(ctx, req); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "register error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "register success, verify code was sent on your email",
	})
}

func (h *AuthHandler) login(c *gin.Context) {
	var (
		req *domain.User
		err error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
		return
	}

	if err = h.authService.Login(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success, verify code was sent on your email",
	})
}

func (h *AuthHandler) verify(c *gin.Context) {
	var (
		err error
		req struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		}
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
		return
	}

	accessToken, err := h.authService.Verify(ctx, req.Email, req.Code)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.Writer.Header().Set("token", accessToken)
	c.JSON(http.StatusOK, gin.H{
		"message": "registration successful",
	})
}
