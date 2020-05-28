package entry

import (
	"net/http"
	"time"
	"go_api_boilerplate/controllers"
	"go_api_boilerplate/db"
	"go_api_boilerplate/middleware"
	"go_api_boilerplate/models"
	"go_api_boilerplate/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func login(c *gin.Context) {
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

func register(c *gin.Context) {
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

func logout(c *gin.Context) {
	// Gets the token metadata from request and handle any possible errors
	token, getErr := middleware.GetTokenData(c.Request)
	if utils.LogError(getErr) != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Deletes the cached access from Redis and handle any possible errors
	delErr := db.DeleteRedisAccess(token.AccessUUID)
	if delErr != nil {
		utils.StaticResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	utils.StaticResponse(c, http.StatusOK, utils.LogoutSuccess)
	return
}
