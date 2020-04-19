package handlers

import (
	"context"
	"net/http"
	"time"
	"vivere_api/models"
	"vivere_api/utils"
	"vivere_api/validators"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// RegisterNewUser expects an input JSON containing the following keys:
// (`email`, `password`)
func RegisterNewUser(c *gin.Context) {
	// Creating an user variable
	var user models.User

	// Trying to bind and validate the request
	bindErr := validators.BindRequest(c, &user)
	valErr := validators.ValidateRequest(c, &user)

	// Handling any pre-database errors
	if utils.HandleError(bindErr, valErr) {
		return
	}

	// Checking if user is valid
	// Note that we use `not` as this block expects `true` when there is no user
	uniqueErr := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&models.User{})
	if !utils.HandleError(uniqueErr) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "User already exists.",
		})
		return
	}

	// Declaring new properties
	user.Token = ""
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Inserting model into collection and checking an insertion error
	_, insertErr := userCollection.InsertOne(context.Background(), user)
	if utils.HandleError(insertErr) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "User could not be added.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "User successfully created.",
	})
	return
}
