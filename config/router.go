package config

import (
	"net/http"
	"vivere_api/handlers"

	"github.com/gin-gonic/gin"
)

// ConfigureRouter expects an router argument to configure
// the desired API routes
func ConfigureRouter(r *gin.Engine) {
	// Login
	r.GET("/login", handlers.LoginHandler)

	// Registration
	r.POST("/register", handlers.RegisterHandler)

	// Users
	r.POST("/user", handlers.AddUser)

	// Non-existing
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "This route is not available.",
		})

		return
	})
}
