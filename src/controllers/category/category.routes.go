package category

import (
	"go_api_boilerplate/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Category-based routes
	category := r.Group("/category")
	{
		category.POST("/", middleware.AuthGuard(), create)
		category.GET("/", list)
		category.GET("/:id", find)
		category.DELETE("/:id", middleware.AuthGuard(), delete)
		category.PATCH("/:id", middleware.AuthGuard(), update)
	}
}
