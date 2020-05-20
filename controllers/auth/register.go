package auth

import (
	"net/http"
	"time"
	"vivere_api/controllers"
	"vivere_api/db"
	"vivere_api/models"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// RegisterNewUser expects an input JSON containing the following keys:
// (`email`, `password`)
func RegisterNewUser(c *gin.Context) {
	// Creates an user variable
	var user models.User

	// Binds and validates the incoming request
	// If it return as false, an error might occurred
	if !controllers.BindAndValidateRequest(c, &user) {
		return
	}

	// Finds a model in collection with the same inputted e-mail
	findErr := db.FindOne(db.UserCollection, bson.M{"email": user.Email}, &models.User{})
	if findErr == nil {
		utils.SendResponse(c, http.StatusInternalServerError, utils.DatabaseAlreadyExists)
		return
	}

	// Hashes and salts the user password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Declares new properties
	user.Password = string(hashedPassword)
	user.Token = ""
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Inserts model into collection
	insertErr := db.InsertOne(db.UserCollection, user)
	if insertErr != nil {
		utils.SendResponse(c, http.StatusInternalServerError, utils.DatabaseInsertionError)
		return
	}

	// Handles a positive response if no errors were found
	utils.SendResponse(c, http.StatusCreated, utils.DatabaseInsertionSuccess)
	return
}
