package handlers

import (
	"context"
	"log"
	"net/http"
	"time"
	"vivere_api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Global variable for the user collection
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

	// Binding the request JSON
	c.BindJSON(&user)

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
