package server

import (
	"net/http"
	"vivere_api/controllers/entry"
	"vivere_api/controllers/item"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter expects an router argument to configure the desired API routes
func InitRouter(r *gin.Engine) {
	// Existing routes
	v1 := r.Group("/v1")
	{
		// Entry-related routes, i.e., login, register and logout
		entry.CreateRoutes(v1)

		// Item-related routes
		item.CreateRoutes(v1)
	}

	// Non-existing routes
	r.NoRoute(func(c *gin.Context) {
		utils.StaticResponse(c, http.StatusNotFound, utils.NoRouteMessage)
		return
	})
}
