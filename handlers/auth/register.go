package auth

import (
	"context"
	"net/http"
	"time"
	"vivere_api/db"
	"vivere_api/handlers"
	"vivere_api/models"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// RegisterNewUser expects an input JSON containing the following keys:
// (`email`, `password`)
func RegisterNewUser(c *gin.Context) {
	// Creating an user variable
	var user models.User

	// Validating the incoming request
	if !handlers.ValidateIncomingRequest(c, &user) {
		return
	}

	// Checking if user is valid
	// Note that we use `not` as this block expects `true` when there is no user
	uniqueErr := db.UserCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&models.User{})
	if utils.HandleError(uniqueErr) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "User already exists.",
		})
		return
	}

	// Hashing and salting the user password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Declaring new properties
	user.Password = string(hashPassword)
	user.Token = ""
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Inserting model into collection and checking an insertion error
	_, insertErr := db.UserCollection.InsertOne(context.Background(), user)
	if !utils.HandleError(insertErr) {
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
