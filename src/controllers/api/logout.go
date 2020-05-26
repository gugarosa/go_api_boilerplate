package api

import (
	"net/http"
	"vivere_api/db"
	"vivere_api/middleware"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
)

// Logout expects no input JSON
func Logout(c *gin.Context) {
	// Gets the token metadata from request and handle any possible errors
	token, getErr := middleware.GetTokenData(c.Request)
	if utils.LogError(getErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Deletes the cached access from Redis and handle any possible errors
	delErr := db.DeleteRedisAccess(token.AccessUUID)
	if delErr != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	utils.StaticResponse(c, http.StatusOK, utils.LogoutSuccess)
	return
}
