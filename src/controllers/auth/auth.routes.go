package auth

import (
	"go_api_boilerplate/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Auth-based routes
	auth := r.Group("/")
	{
		auth.POST("/login", login)
		auth.POST("/register", register)
		auth.POST("/logout", middleware.AuthGuard(), logout)
	}
}
