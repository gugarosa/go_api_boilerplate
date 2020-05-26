package entry

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

// Login expects an input JSON containing the following keys:
// (`email`, `password`)
func Login(c *gin.Context) {
	// Creates input and database user variables
	var inputUser, dbUser models.User

	// Binds and validates the model, and handles any possible errors
	checkErr := controllers.BindAndValidateRequest(c, &inputUser)
	if utils.LogError(checkErr) != nil {
		utils.StaticResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Finds a model in collection with the same inputted e-mail
	findErr := db.FindOne(db.UserCollection, bson.M{"email": inputUser.Email}, &dbUser)
	if utils.LogError(findErr) != nil {
		utils.StaticResponse(c, http.StatusInternalServerError, utils.LoginError)
		return
	}

	// Compares if both passwords are the same
	passwordErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password))
	if utils.LogError(passwordErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	// Creates new authentication tokens
	token, tokenErr := middleware.CreateToken(dbUser.ID)
	if utils.LogError(tokenErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	// Sets the cached accesses in Redis
	redisErr := db.CreateRedisAccess(dbUser.ID, token)
	if utils.LogError(redisErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
	return
}
