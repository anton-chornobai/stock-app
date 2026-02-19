package handlers

import (
	"auth-service/internal/application"
	"auth-service/internal/domain"
	"errors"
	"log/slog"
	"net/http"

	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	authService application.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService application.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	token, err := h.authService.Signup(ctx, req.Email, req.Password)

	if err != nil {
		h.logger.Warn("signup error", "err", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(201, gin.H{
		"token": token,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "couldnt parse credentials"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	token, err := h.authService.Login(ctx, req.Email, req.Password)

	if err != nil {

		if errors.Is(err, domain.ErrInvalidEmail) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "couldnt login user"})

		}
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
