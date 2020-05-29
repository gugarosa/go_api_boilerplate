package server

import (
	"go_api_boilerplate/controllers/auth"
	"go_api_boilerplate/controllers/item"
	"go_api_boilerplate/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter expects an router argument to configure the desired API routes
func InitRouter(r *gin.Engine) {
	// Existing routes
	v1 := r.Group("/v1")
	{
		// Auth-related routes, i.e., login, register and logout
		auth.CreateRoutes(v1)

		// Item-related routes
		item.CreateRoutes(v1)
	}

	// Non-existing routes
	r.NoRoute(func(c *gin.Context) {
		utils.ConstantResponse(c, http.StatusNotFound, utils.NoRouteMessage)
		return
	})
}
