package auth

import (
	"time"
	"vivere_api/db"
	"vivere_api/handlers"
	"vivere_api/models"

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
	// If it return as false, an error might occurred
	if !handlers.ValidateIncomingRequest(c, &user) {
		return
	}

	// Finding a model in collection with the same inputted e-mail
	if !db.FindOne(c, db.UserCollection, bson.M{"email": user.Email}, &models.User{}) {
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
