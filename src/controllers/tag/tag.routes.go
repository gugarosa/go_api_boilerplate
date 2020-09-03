package tag

import (
	"go_api_boilerplate/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Tag-based routes
	tag := r.Group("/tag")
	{
		tag.POST("/", middleware.AuthGuard(), create)
		tag.GET("/", list)
		tag.GET("/:id", find)
		tag.DELETE("/:id", middleware.AuthGuard(), delete)
		tag.PATCH("/:id", middleware.AuthGuard(), update)
	}
}
