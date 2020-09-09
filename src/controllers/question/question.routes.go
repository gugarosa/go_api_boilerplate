package question

import (
	"go_api_boilerplate/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoutes expects a RouterGroup
// to create a group of common-knowledge routes
func CreateRoutes(r *gin.RouterGroup) {
	// Question-based routes
	question := r.Group("/question")
	{
		question.POST("/", middleware.AuthGuard(), create)
		question.GET("/", list)
		question.GET("/:id", find)
		question.DELETE("/:id", middleware.AuthGuard(), delete)
		question.PATCH("/:id", middleware.AuthGuard(), update)
	}
}
