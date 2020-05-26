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
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := db.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}
