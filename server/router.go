package server

import (
	"net/http"
	"vivere_api/controllers/auth"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter expects an router argument to configure the desired API routes
func InitRouter(r *gin.Engine) {
	// Authorization
	r.POST("/login", auth.LogNewUser)
	r.POST("/register", auth.RegisterNewUser)

	// Non-existing
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": utils.NoRouteMessage,
		})

		return
	})
}
