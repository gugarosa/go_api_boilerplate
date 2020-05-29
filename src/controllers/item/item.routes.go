package item

import (
	"go_api_boilerplate/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Item-based routes
	item := r.Group("/item")
	{
		item.POST("/", middleware.AuthGuard(), create)
		item.GET("/", list)
		item.GET("/:id", find)
		item.DELETE("/:id", middleware.AuthGuard(), delete)
		item.PATCH("/:id", middleware.AuthGuard(), update)
	}
}
