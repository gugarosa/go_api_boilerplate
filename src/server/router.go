package server

import (
	"go_api_boilerplate/controllers/auth"
	"go_api_boilerplate/controllers/product"
	"go_api_boilerplate/controllers/survey"
	"go_api_boilerplate/controllers/tag"
	"go_api_boilerplate/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter expects an router argument to configure the desired API routes
func InitRouter(r *gin.Engine) {
	// Existing routes
	v1 := r.Group("/v1")
	{
		// Auth-related routes, i.e., login, refresh, register and logout
		auth.CreateRoutes(v1)

		// Product-related routes
		product.CreateRoutes(v1)

		// Survey-related routes
		survey.CreateRoutes(v1)

		// Tag-related routes
		tag.CreateRoutes(v1)
	}

	// Non-existing routes
	r.NoRoute(func(c *gin.Context) {
		utils.ConstantResponse(c, http.StatusNotFound, utils.NoRouteMessage)
		return
	})
}
