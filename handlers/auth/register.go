package auth

import (
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

	// Binding and validating the incoming request
	// If it return as false, an error might occurred
	if !handlers.BindAndValidateIncomingRequest(c, &user) {
		return
	}

	// Finding a model in collection with the same inputted e-mail
	if db.FindOne(c, db.UserCollection, bson.M{"email": user.Email}, &models.User{}) {
		// If a model has been found, return a JSON indicating that user already exists
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": utils.UserAlreadyExists,
		})
		return
	}

	// Hashing and salting the user password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Declaring new properties
	user.Password = string(hashedPassword)
	user.Token = ""
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Inserting model into collection and checking for an insertion error
	if !db.InsertOne(c, db.UserCollection, user) {
		return
	}

	return
}
