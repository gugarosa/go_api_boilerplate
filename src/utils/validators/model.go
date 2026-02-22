package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Reusable validator instance
var validate = validator.New()

// BindModel expects a context and a model
// to bind a generic model to the required context
func BindModel(c *gin.Context, model interface{}) error {
	err := c.ShouldBindJSON(model)
	if err != nil {
		return err
	}

	return nil
}

// ValidateModel expects a model
// to validate it based on pre-defined validation rules
func ValidateModel(model interface{}) error {
	err := validate.Struct(model)
	if err != nil {
		return err
	}

	return nil
}
