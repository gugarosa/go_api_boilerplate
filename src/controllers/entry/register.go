package entry

import (
	"net/http"
	"time"
	"vivere_api/db"
	"vivere_api/models"
	"vivere_api/utils"
	"vivere_api/utils/validators"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Register expects an input JSON containing the following keys:
// (`email`, `password`)
func Register(c *gin.Context) {
	// Creates an user variable
	var user models.User

	// Binds the model and handle any possible errors
	bindErr := validators.BindModel(c, &user)
	if utils.HandleError(bindErr) != nil {
		return
	}

	// Validates the model and handle any possible errors
	valErr := validators.ValidateModel(c, &user)
	if utils.HandleError(valErr) != nil {
		return
	}

	// Finds a model in collection with the same inputted e-mail
	findErr := db.FindOne(db.UserCollection, bson.M{"email": user.Email}, &models.User{})
	if utils.HandleError(findErr) == nil {
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
	if utils.HandleError(insertErr) != nil {
		utils.StaticResponse(c, http.StatusInternalServerError, utils.DatabaseInsertionError)
		return
	}

	utils.StaticResponse(c, http.StatusCreated, utils.DatabaseInsertionSuccess)
	return
}
