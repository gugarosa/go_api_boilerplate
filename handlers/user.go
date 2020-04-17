package handlers

import (
	"context"
	"log"
	"net/http"
	"time"
	"vivere_api/models"
	"vivere_api/validators"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Global variables for the UserHandler
var collection *mongo.Collection

// UserCollection expects a MongoDB database as parameter and sets in the scope
// a variable to the `users` collection
func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}

// AddUser expects an input JSON containing the following keys:
// (`email`, `password`)
func AddUser(c *gin.Context) {
	// Creating an user variable
	var user models.User

	// Trying to bind the request
	berr := validators.BindRequest(c, &user)

	// Trying to validate the request
	verr := validators.ValidateRequest(c, &user)

	// Checking if there is any error on request binding or validation
	if berr != nil || verr != nil {
		log.Println(berr, verr)

		return
	}

	// Declaring new properties
	user.Token = ""
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Inserting model into collection
	_, err := collection.InsertOne(context.Background(), user)

	// Checking if there is an error and returning error status and message
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "User could not be added to collection.",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "User successfully created.",
	})

	return
}
