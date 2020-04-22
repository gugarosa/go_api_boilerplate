package db

import (
	"context"
	"net/http"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertOne ...
func InsertOne(c *gin.Context, collection *mongo.Collection, model interface{}) {
	_, insertErr := db.UserCollection.InsertOne(context.Background(), user)
	if !utils.HandleError(insertErr) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "User could not be added.",
		})
		return
	}
}
