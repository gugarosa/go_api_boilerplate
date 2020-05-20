package auth

import (
	"net/http"
	"vivere_api/controllers"
	"vivere_api/db"
	"vivere_api/middleware"
	"vivere_api/models"
	"vivere_api/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// LogUser expects an input JSON containing the following keys:
// (`email`, `password`)
func LogUser(c *gin.Context) {
	// Creates input and database user variables
	var inputUser, dbUser models.User

	// Binds the request
	bindErr := controllers.BindRequest(c, &inputUser)
	if bindErr != nil {
		return
	}

	// Finds a model in collection with the same inputted e-mail
	findErr := db.FindOne(db.UserCollection, bson.M{"email": inputUser.Email}, &dbUser)
	if findErr != nil {
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
	token, tokenErr := middleware.CreateNewToken(dbUser.ID)
	if utils.HandleError(tokenErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	// Sets the tokens in Redis
	setErr := db.SetTokens(dbUser.ID, token)
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
