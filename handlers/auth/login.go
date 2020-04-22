package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginHandler handles the login request
func LoginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
