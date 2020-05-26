package validators

import (
	"net/http"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindModel expects a context and a model
// to bind a generic model to the required context
func BindModel(c *gin.Context, model interface{}) error {
	err := c.ShouldBind(&model)
	if err != nil {
		utils.StaticResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

// ValidateModel expects a context and a model
// to validate a request based on pre-defined validation rules
func ValidateModel(c *gin.Context, model interface{}) error {
	err := validator.New().Struct(model)
	if err != nil {
		utils.StaticResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}
