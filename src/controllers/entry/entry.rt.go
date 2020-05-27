package entry

import (
	"vivere_api/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Entry-based routes
	entry := r.Group("/")
	{
		entry.POST("/login", login)
		entry.POST("/register", register)
		entry.POST("/logout", middleware.AuthGuard(), logout)
	}
}
