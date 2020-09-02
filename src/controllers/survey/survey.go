package survey

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
	// Creates an empty Survey model variable
	var survey models.Survey

	// Authenticates the request and handle any possible errors
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &survey)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Declares new properties
	survey.Active = true
	survey.CreatedAt = time.Now()
	survey.UpdatedAt = time.Now()

	// Inserts model into collection
	insertErr := database.InsertOne(database.SurveyCollection, survey)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
	return
}

func list(c *gin.Context) {
	// Finds a model in collection with the same inputted ID
	surveys, findErr := database.FindAll(database.SurveyCollection)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": surveys,
	})
	return

}

func find(c *gin.Context) {
	// Creates variable to hold the found model
	var survey bson.M

	// Gathers the string ID as an ObjectID
	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestParameterError)
		return
	}

	// Finds a model in collection with the same inputted ID
	findErr := database.FindOne(database.SurveyCollection, bson.M{"_id": id}, &survey)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": survey,
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
	deleteErr := database.DeleteOne(database.SurveyCollection, id)
	if utils.LogError(deleteErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseDeleteError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseDeleteSuccess)
	return
}

func update(c *gin.Context) {
	// Creates an empty Survey model variable
	var survey models.Survey

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
	checkErr := controllers.BindAndValidateRequest(c, &survey)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Declares new properties
	survey.UpdatedAt = time.Now()

	// Decodes the structure into an update document
	update, decodeErr := controllers.DecodeStruct(survey)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	// Updates a model in collection using filter and request data
	updateErr := database.UpdateOne(database.SurveyCollection, id, update)
	if utils.LogError(updateErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseUpdateError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseUpdateSuccess)
	return
}
