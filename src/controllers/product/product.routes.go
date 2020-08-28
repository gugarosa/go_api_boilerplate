package product

import (
	"go_api_boilerplate/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Product-based routes
	product := r.Group("/product")
	{
		product.POST("/", middleware.AuthGuard(), create)
		product.GET("/", list)
		product.GET("/:id", find)
		product.DELETE("/:id", middleware.AuthGuard(), delete)
		product.PATCH("/:id", middleware.AuthGuard(), update)
	}
}
