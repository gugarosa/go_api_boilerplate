package api

import (
	"net/http"
	"vivere_api/controllers"
	"vivere_api/models"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
)

// CreateProduct expects an input JSON containing the following keys:
// (`name`)
func CreateProduct(c *gin.Context) {
	// Creates an empty Product model variable
	var product models.Product

	// Authenticates the request and handle any possible errors
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.StaticResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &product)
	if utils.LogError(checkErr) != nil {
		return
	}

	//
	product.Name = "Test"

	utils.StaticResponse(c, http.StatusCreated, product.Name)
	return
}
