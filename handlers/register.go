package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles an user's registration
func RegisterHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
