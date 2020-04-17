package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindRequest tries to bind a generic model to the required context
func BindRequest(c *gin.Context, model interface{}) error {
	if err := c.ShouldBind(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})

		return err
	}

	return nil
}

// ValidateRequest tries to validate a request based on pre-defined validation rules
func ValidateRequest(c *gin.Context, model interface{}) error {
	if err := validator.New().Struct(model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})

		return err
	}

	return nil
}
