package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginHandle handles the login request
func LoginHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
