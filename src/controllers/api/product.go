package api

import (
	"net/http"
	"vivere_api/db"
	"vivere_api/middleware"
	"vivere_api/models"

	"github.com/gin-gonic/gin"
)

// CreateProduct ...
func CreateProduct(c *gin.Context) {
	var td *models.Product
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	tokenAuth, err := middleware.GetTokenData(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	_, err = db.GetAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	td.Name = "Test"

	//you can proceed to save the Product to a database
	//but we will just return it to the caller here:
	c.JSON(http.StatusCreated, td)
}
