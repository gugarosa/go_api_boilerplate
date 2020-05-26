package server

import (
	"net/http"
	"vivere_api/controllers/api"
	"vivere_api/controllers/entry"
	"vivere_api/middleware"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter expects an router argument to configure the desired API routes
func InitRouter(r *gin.Engine) {
	// Non-authenticated routes
	r.POST("/login", entry.Login)
	r.POST("/register", entry.Register)

	// Authenticated routes
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthGuard())
	{
		apiGroup.POST("/product", api.CreateProduct)
		apiGroup.POST("/logout", api.Logout)
	}

	// Non-existing routes
	r.NoRoute(func(c *gin.Context) {
		utils.StaticResponse(c, http.StatusNotFound, utils.NoRouteMessage)
		return
	})
}
