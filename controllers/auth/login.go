package auth

import (
	"net/http"
	"vivere_api/db"
	"vivere_api/controllers"
	"vivere_api/middleware"
	"vivere_api/models"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// LogNewUser expects an input JSON containing the following keys:
// (`email`, `password`)
func LogNewUser(c *gin.Context) {
	// Creating an input and database user variables
	var inputUser, dbUser models.User

	// Binding the incoming request
	// If it return as false, an error might occurred
	if !controllers.BindIncomingRequest(c, &inputUser) {
		return
	}

	// Finding a model in collection with the same inputted e-mail
	if !db.FindOne(c, db.UserCollection, bson.M{"email": inputUser.Email}, &dbUser) {
		// If a model has not been found, return a JSON indicating that user does not exists
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": utils.LoginError,
		})
		return
	}

	// Comparing if both passwords are the same
	passwordErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password))

	// Handles if an error has occured
	if !utils.HandleError(passwordErr) {
		// If yes, returns a JSON with an error status
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": utils.LoginError,
		})
		return
	}

	//
	token, tokenErr := middleware.CreateNewToken(dbUser.ID)

	// Handles if an error has occured
	if !utils.HandleError(tokenErr) {
		// If yes, returns a JSON with an error status
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": utils.LoginError,
		})
		return
	}

	saveErr := middleware.CreateAuth(dbUser.ID, token)

	// Handles if an error has occured
	if !utils.HandleError(saveErr) {
		// If yes, returns a JSON with an error status
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": utils.LoginError,
		})
		return
	}

	// If not, returns a JSON with a success status
	c.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}
