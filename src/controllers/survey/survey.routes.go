package survey

import (
	"go_api_boilerplate/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Survey-based routes
	survey := r.Group("/survey")
	{
		survey.POST("/", middleware.AuthGuard(), create)
		survey.GET("/", list)
		survey.GET("/:id", find)
		survey.DELETE("/:id", middleware.AuthGuard(), delete)
		survey.PATCH("/:id", middleware.AuthGuard(), update)
	}
}
