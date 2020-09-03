package tag

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
	// Creates an empty Tag model variable
	var tag models.Tag

	// Authenticates the request and handle any possible errors
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &tag)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Declares new properties
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	// Inserts model into collection
	insertErr := database.InsertOne(database.TagCollection, tag)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
	return
}

func list(c *gin.Context) {
	// Finds a model in collection with the same inputted ID
	tags, findErr := database.FindAll(database.TagCollection)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": tags,
	})
	return

}

func find(c *gin.Context) {
	// Creates variable to hold the found model
	var tag bson.M

	// Gathers the string ID as an ObjectID
	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestParameterError)
		return
	}

	// Finds a model in collection with the same inputted ID
	findErr := database.FindOne(database.TagCollection, bson.M{"_id": id}, &tag)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": tag,
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
	deleteErr := database.DeleteOne(database.TagCollection, id)
	if utils.LogError(deleteErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseDeleteError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseDeleteSuccess)
	return
}

func update(c *gin.Context) {
	// Creates an empty Tag model variable
	var tag models.Tag

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
	checkErr := controllers.BindAndValidateRequest(c, &tag)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Declares new properties
	tag.UpdatedAt = time.Now()

	// Decodes the structure into an update document
	update, decodeErr := controllers.DecodeStruct(tag)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	// Updates a model in collection using filter and request data
	updateErr := database.UpdateOne(database.TagCollection, id, update)
	if utils.LogError(updateErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseUpdateError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseUpdateSuccess)
	return
}
