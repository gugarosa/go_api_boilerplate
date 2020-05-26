package controllers

import (
	"vivere_api/db"
	"vivere_api/middleware"
	"vivere_api/utils/validators"

	"github.com/gin-gonic/gin"
)

// AuthRequest expects a context
// to authenticate the request
func AuthRequest(c *gin.Context) error {
	// Gets the token metadata from request and handle any possible errors
	token, err := middleware.GetTokenData(c.Request)
	if err != nil {
		return err
	}

	// Gets the cached access from Redis and handle any possible errors
	err = db.GetRedisAccess(token)
	if err != nil {
		return err
	}

	return nil
}

// BindAndValidateRequest expects a context and a model
// to bind and validate the request with the model
func BindAndValidateRequest(c *gin.Context, model interface{}) error {
	// Binds the request to the model
	err := validators.BindModel(c, &model)
	if err != nil {
		return err
	}

	// Validates the model
	err = validators.ValidateModel(model)
	if err != nil {
		return err
	}

	return nil
}
