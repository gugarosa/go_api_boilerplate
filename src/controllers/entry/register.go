package entry

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

// Register expects an input JSON containing the following keys:
// (`email`, `password`)
func Register(c *gin.Context) {
	// Creates an user variable
	var user models.User

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &user)
	if utils.LogError(checkErr) != nil {
		utils.StaticResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Finds a model in collection with the same inputted e-mail
	findErr := db.FindOne(db.UserCollection, bson.M{"email": user.Email}, &models.User{})
	if utils.LogError(findErr) == nil {
		utils.StaticResponse(c, http.StatusInternalServerError, utils.DatabaseAlreadyExists)
		return
	}

	// Hashes and salts the user password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Declares new properties
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Inserts model into collection
	insertErr := db.InsertOne(db.UserCollection, user)
	if utils.LogError(insertErr) != nil {
		utils.StaticResponse(c, http.StatusInternalServerError, utils.DatabaseInsertionError)
		return
	}

	utils.StaticResponse(c, http.StatusCreated, utils.DatabaseInsertionSuccess)
	return
}
