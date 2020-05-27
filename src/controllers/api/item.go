package api

import (
	"net/http"
	"time"
	"vivere_api/controllers"
	"vivere_api/db"
	"vivere_api/models"
	"vivere_api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// CreateItem expects an input JSON containing the following keys:
// (`name`)
func CreateItem(c *gin.Context) {
	// Creates an empty Item model variable
	var item models.Item

	// Authenticates the request and handle any possible errors
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &item)
	if utils.LogError(checkErr) != nil {
		utils.StaticResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Declares new properties
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	// Inserts model into collection
	insertErr := db.InsertOne(db.ItemCollection, item)
	if utils.LogError(insertErr) != nil {
		utils.StaticResponse(c, http.StatusInternalServerError, utils.DatabaseInsertionError)
		return
	}

	utils.StaticResponse(c, http.StatusCreated, utils.DatabaseInsertionSuccess)
	return
}

// GetItemByID ...
func GetItemByID(c *gin.Context) {
	// Creates item variable
	var item bson.M

	//
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	// Finds a model in collection with the same inputted e-mail
	findErr := db.FindOne(db.ItemCollection, bson.M{"_id": id}, &item)
	if utils.LogError(findErr) != nil {
		utils.StaticResponse(c, http.StatusInternalServerError, utils.DatabaseNonExists)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, item)
	return

}
