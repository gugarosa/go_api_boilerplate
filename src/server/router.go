package server

import (
	"net/http"
	"vivere_api/controllers/api"
	"vivere_api/controllers/auth"
	"vivere_api/middleware"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter expects an router argument to configure the desired API routes
func InitRouter(r *gin.Engine) {
	// Authorization
	r.POST("/login", auth.Login)
	r.POST("/register", auth.RegisterUser)
	r.POST("/logout", middleware.TokenAuthMiddleware(), auth.Logout)

	// Product
	r.POST("/product", middleware.TokenAuthMiddleware(), api.CreateProduct)

	// Non-existing
	r.NoRoute(func(c *gin.Context) {
		utils.StaticResponse(c, http.StatusNotFound, utils.NoRouteMessage)
		return
	})
}
