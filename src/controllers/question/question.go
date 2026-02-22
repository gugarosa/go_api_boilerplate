package question

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
	var question models.Question

	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	checkErr := controllers.BindAndValidateRequest(c, &question)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	question.Active = true
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	insertErr := database.InsertOne(c.Request.Context(), database.QuestionCollection, question)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
}

func list(c *gin.Context) {
	pipeline := []bson.M{
		bson.M{"$lookup": bson.M{"from": "tags", "localField": "tags", "foreignField": "_id", "as": "tags"}},
	}

	questions, findErr := database.FindAllWithAggregate(c.Request.Context(), database.QuestionCollection, pipeline)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": questions,
	})
}

func find(c *gin.Context) {
	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestParameterError)
		return
	}

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"_id": id}},
		bson.M{"$lookup": bson.M{"from": "tags", "localField": "tags", "foreignField": "_id", "as": "tags"}},
	}

	question, findErr := database.FindOneWithAggregate(c.Request.Context(), database.QuestionCollection, pipeline)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": question,
	})
}

func delete(c *gin.Context) {
	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestParameterError)
		return
	}

	deleteErr := database.DeleteOne(c.Request.Context(), database.QuestionCollection, id)
	if utils.LogError(deleteErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseDeleteError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseDeleteSuccess)
}

func update(c *gin.Context) {
	var question models.Question

	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestParameterError)
		return
	}

	checkErr := controllers.BindAndValidateRequest(c, &question)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	question.UpdatedAt = time.Now()

	update, decodeErr := controllers.DecodeStruct(question)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	updateErr := database.UpdateOne(c.Request.Context(), database.QuestionCollection, id, update)
	if utils.LogError(updateErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseUpdateError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseUpdateSuccess)
}
