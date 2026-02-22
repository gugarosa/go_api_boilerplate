package controllers

import (
	"errors"
	"go_api_boilerplate/database"
	"go_api_boilerplate/models"
	"go_api_boilerplate/utils/validators"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// AuthRequest expects a context
// to authenticate the request via cached Redis session
func AuthRequest(c *gin.Context) error {
	// Retrieve token data stored by AuthGuard middleware
	val, exists := c.Get("token_data")
	if !exists {
		return errors.New("missing token data")
	}

	token, ok := val.(*models.RedisAccess)
	if !ok {
		return errors.New("invalid token data")
	}

	// Gets the cached access from Redis and handle any possible errors
	err := database.GetRedisAccess(c.Request.Context(), token)
	if err != nil {
		return err
	}

	return nil
}

// BindAndValidateRequest expects a context and a model
// to bind and validate the request with the model
func BindAndValidateRequest(c *gin.Context, model interface{}) error {
	// Binds the request to the model
	err := validators.BindModel(c, model)
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

// DecodeStruct expects an interface (struct)
// to decode it into a BSON document
func DecodeStruct(s interface{}) (bson.M, error) {
	// Creating a bson.M variable
	var decoded bson.M

	// Marshalling the input struct
	encoded, err := bson.Marshal(s)
	if err != nil {
		return nil, err
	}

	// Unmarshalling the encoded object
	err = bson.Unmarshal(encoded, &decoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

// EncodeStruct expects a BSON document and an interface (struct)
// to encode it back to the struct
func EncodeStruct(m bson.M, encoded interface{}) error {
	// Marshalling the input BSON
	decoded, err := bson.Marshal(m)
	if err != nil {
		return err
	}

	// Unmarshalling the encoded object
	err = bson.Unmarshal(decoded, encoded)
	if err != nil {
		return err
	}

	return nil
}
