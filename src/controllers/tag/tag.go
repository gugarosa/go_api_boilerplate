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
	var tag models.Tag

	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	checkErr := controllers.BindAndValidateRequest(c, &tag)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	insertErr := database.InsertOne(c.Request.Context(), database.TagCollection, tag)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
}

func list(c *gin.Context) {
	tags, findErr := database.FindAll(c.Request.Context(), database.TagCollection)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": tags,
	})
}

func find(c *gin.Context) {
	id, hexErr := primitive.ObjectIDFromHex(c.Param("id"))
	if utils.LogError(hexErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestParameterError)
		return
	}

	tag, findErr := database.FindOne(c.Request.Context(), database.TagCollection, bson.M{"_id": id})
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": tag,
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

	deleteErr := database.DeleteOne(c.Request.Context(), database.TagCollection, id)
	if utils.LogError(deleteErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseDeleteError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseDeleteSuccess)
}

func update(c *gin.Context) {
	var tag models.Tag

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

	checkErr := controllers.BindAndValidateRequest(c, &tag)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	tag.UpdatedAt = time.Now()

	update, decodeErr := controllers.DecodeStruct(tag)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	updateErr := database.UpdateOne(c.Request.Context(), database.TagCollection, id, update)
	if utils.LogError(updateErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseUpdateError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseUpdateSuccess)
}
