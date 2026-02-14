package handlers

import "github.com/gin-gonic/gin"

type Handler struct {}

func (h *Handler) Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "login",
	})
}