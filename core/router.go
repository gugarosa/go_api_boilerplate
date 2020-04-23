package core

import (
	"net/http"
	"vivere_api/handlers/auth"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
)

// AddRouter expects an router argument to configure
// the desired API routes
func AddRouter(r *gin.Engine) {
	// Login
	r.POST("/login", auth.LogNewUser)

	// Registration
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
