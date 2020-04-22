package handlers

import (
	"vivere_api/utils"
	"vivere_api/utils/validators"

	"github.com/gin-gonic/gin"
)

// ValidateIncomingRequest tries to bind and validate the request,
// and returns a boolean whether any error has happened or not
func ValidateIncomingRequest(c *gin.Context, model interface{}) bool {
	// Trying to bind the request
	bindErr := validators.BindRequest(c, &model)

	// Trying to validate the request
	valErr := validators.ValidateRequest(c, model)

	// Handling and returning any possible errors
	return utils.HandleError(bindErr, valErr)
}
