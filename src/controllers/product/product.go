package product

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
	var product models.Product

	authErr := controllers.AuthRequest(c)
	if utils.LogError(authErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	checkErr := controllers.BindAndValidateRequest(c, &product)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	product.Active = true
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	insertErr := database.InsertOne(c.Request.Context(), database.ProductCollection, product)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
}

func list(c *gin.Context) {
	pipeline := []bson.M{
		bson.M{"$lookup": bson.M{"from": "tags", "localField": "tags", "foreignField": "_id", "as": "tags"}},
		bson.M{"$lookup": bson.M{"from": "categories", "localField": "categories", "foreignField": "_id", "as": "categories"}},
	}

	products, findErr := database.FindAllWithAggregate(c.Request.Context(), database.ProductCollection, pipeline)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": products,
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
		bson.M{"$lookup": bson.M{"from": "categories", "localField": "categories", "foreignField": "_id", "as": "categories"}},
	}

	product, findErr := database.FindOneWithAggregate(c.Request.Context(), database.ProductCollection, pipeline)
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseFindError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"response": product,
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

	deleteErr := database.DeleteOne(c.Request.Context(), database.ProductCollection, id)
	if utils.LogError(deleteErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseDeleteError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseDeleteSuccess)
}

func update(c *gin.Context) {
	var product models.Product

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

	checkErr := controllers.BindAndValidateRequest(c, &product)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	product.UpdatedAt = time.Now()

	update, decodeErr := controllers.DecodeStruct(product)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	updateErr := database.UpdateOne(c.Request.Context(), database.ProductCollection, id, update)
	if utils.LogError(updateErr) != nil {
		utils.ConstantResponse(c, http.StatusNotFound, utils.DatabaseUpdateError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.DatabaseUpdateSuccess)
}
