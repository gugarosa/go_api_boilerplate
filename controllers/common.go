package controllers

import (
	"vivere_api/utils"
	"vivere_api/utils/validators"

	"github.com/gin-gonic/gin"
)

// BindRequest tries to bind the request,
// and returns a whether any error has happened or not
func BindRequest(c *gin.Context, model interface{}) error {
	// Trying to bind the request
	bindErr := validators.BindRequest(c, &model)

	return utils.HandleError(bindErr)
}

// BindAndValidateRequest tries to bind and validate the request,
// and returns a boolean whether any error has happened or not
func BindAndValidateRequest(c *gin.Context, model interface{}) error {
	// Tries to bind the request
	bindErr := validators.BindRequest(c, &model)

	// Tries to validate the request
	valErr := validators.ValidateRequest(c, model)

	return utils.HandleError(bindErr, valErr)
}
