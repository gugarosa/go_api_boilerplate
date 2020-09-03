package category

import (
	"go_api_boilerplate/controllers"
	"go_api_boilerplate/database"
	"go_api_boilerplate/models"
	"go_api_boilerplate/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

func create(c *gin.Context) {
	// Creates an empty Category model variable
	var category models.Category

	// Authenticates the request and handle any possible errors
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &category)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Declares new properties
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	// Inserts model into collection
	insertErr := database.InsertOne(database.CategoryCollection, category)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
	return
}

func list(c *gin.Context) {
	// Finds a model in collection with the same inputted ID
	categories, findErr := database.FindAll(database.CategoryCollection)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": categories,
	})
	return

}

func find(c *gin.Context) {
	// Creates variable to hold the found model
	var category bson.M

	// Gathers the string ID as an ObjectID
	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestParameterError)
		return
	}

	// Finds a model in collection with the same inputted ID
	findErr := database.FindOne(database.CategoryCollection, bson.M{"_id": id}, &category)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": category,
	})
	return

}

func delete(c *gin.Context) {
	// Authenticates the request and handle any possible errors
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Gathers the string ID as an ObjectID
	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestParameterError)
		return
	}

	// Deletes a model in collection with the same inputted ID
	deleteErr := database.DeleteOne(database.CategoryCollection, id)
	if utils.LogError(deleteErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseDeleteError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseDeleteSuccess)
	return
}

func update(c *gin.Context) {
	// Creates an empty Category model variable
	var category models.Category

	// Authenticates the request and handle any possible errors
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Gathers the string ID as an ObjectID
	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestParameterError)
		return
	}

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &category)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Declares new properties
	category.UpdatedAt = time.Now()

	// Decodes the structure into an update document
	update, decodeErr := controllers.DecodeStruct(category)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	// Updates a model in collection using filter and request data
	updateErr := database.UpdateOne(database.CategoryCollection, id, update)
	if utils.LogError(updateErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseUpdateError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseUpdateSuccess)
	return
}
