package entry

import (
	"net/http"
	"vivere_api/db"
	"vivere_api/middleware"
	"vivere_api/models"
	"vivere_api/utils"
	"vivere_api/utils/validators"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Login expects an input JSON containing the following keys:
// (`email`, `password`)
func Login(c *gin.Context) {
	// Creates input and database user variables
	var inputUser, dbUser models.User

	// Binds the request and handle any possible errors
	bindErr := validators.BindModel(c, &inputUser)
	if utils.HandleError(bindErr) != nil {
		return
	}

	// Finds a model in collection with the same inputted e-mail
	findErr := db.FindOne(db.UserCollection, bson.M{"email": inputUser.Email}, &dbUser)
	if utils.HandleError(findErr) != nil {
		utils.StaticResponse(c, http.StatusInternalServerError, utils.LoginError)
		return
	}

	// Compares if both passwords are the same
	passwordErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password))
	if utils.HandleError(passwordErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	// Creates new authentication tokens
	token, tokenErr := middleware.CreateToken(dbUser.ID)
	if utils.HandleError(tokenErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	// Sets the tokens in Redis
	setErr := db.SetAuth(dbUser.ID, token)
	if utils.HandleError(setErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
	return
}
