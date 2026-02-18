package handlers

import (
	"auth-service/internal/application"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService application.AuthService
}

func NewAuthHandler(authService application.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req application.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	token, err := h.authService.Signup(ctx, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})
}
