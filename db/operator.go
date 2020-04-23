package db

import (
	"context"
	"net/http"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOne expects a collection and a model in order to find
// a single document into the database
func FindOne(c *gin.Context, collection *mongo.Collection, filter bson.M, model interface{}) bool {
	// Tries to find model in the database
	uniqueErr := collection.FindOne(c, filter).Decode(model)

	// Handles if an error has occured
	// Note that we use the comparison as `true` as this block expects `false` when there is no model
	if utils.HandleError(uniqueErr) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": utils.DatabaseAlreadyExists,
		})
		return false
	}

	return true
}

// InsertOne expects a collection and a model in order to insert
// a new document into the database
func InsertOne(c *gin.Context, collection *mongo.Collection, model interface{}) bool {
	// Tries to insert the model in the database
	_, insertErr := collection.InsertOne(context.Background(), model)

	// Handles if an error has occured
	if !utils.HandleError(insertErr) {
		// If yes, returns a JSON with an error status
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": utils.DatabaseInsertionError,
		})
		return false
	}

	// If not, returns a JSON with a success status
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": utils.DatabaseInsertionSuccess,
	})
	return true
}
