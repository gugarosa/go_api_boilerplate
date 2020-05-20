package controllers

import (
	"vivere_api/utils"
	"vivere_api/utils/validators"

	"github.com/gin-gonic/gin"
)

// BindIncomingRequest tries to bind the request,
// and returns a boolean whether any error has happened or not
func BindIncomingRequest(c *gin.Context, model interface{}) bool {
	// Trying to bind the request
	bindErr := validators.BindRequest(c, &model)

	// Handling and returning any possible errors
	return utils.HandleError(bindErr)
}

// BindAndValidateIncomingRequest tries to bind and validate the request,
// and returns a boolean whether any error has happened or not
func BindAndValidateIncomingRequest(c *gin.Context, model interface{}) bool {
	// Trying to bind the request
	bindErr := validators.BindRequest(c, &model)

	// Trying to validate the request
	valErr := validators.ValidateRequest(c, model)

	// Handling and returning any possible errors
	return utils.HandleError(bindErr, valErr)
}
