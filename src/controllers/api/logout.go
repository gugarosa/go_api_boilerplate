package api

import (
	"net/http"
	"vivere_api/db"
	"vivere_api/middleware"

	"github.com/gin-gonic/gin"
)

// Logout ...
func Logout(c *gin.Context) {
	au, err := middleware.GetTokenData(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized1")
		return
	}
	delErr := db.DeleteRedisAccess(au.AccessUUID)
	if delErr != nil { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized2")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}
