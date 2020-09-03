package auth

import (
	"go_api_boilerplate/controllers"
	"go_api_boilerplate/database"
	"go_api_boilerplate/middleware"
	"go_api_boilerplate/models"
	"go_api_boilerplate/utils"
	"net/http"
	"time"

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
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Finds a model in collection with the same inputted e-mail
	user, findErr := database.FindOne(database.UserCollection, bson.M{"email": inputUser.Email})
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.LoginError)
		return
	}

	// Encodes the document into a structure
	decodeErr := controllers.EncodeStruct(user, &dbUser)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	// Compares if both passwords are the same
	passwordErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password))
	if utils.LogError(passwordErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	// Creates new authentication tokens
	token, tokenErr := middleware.CreateToken(dbUser.ID)
	if utils.LogError(tokenErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	// Sets the cached accesses in Redis
	redisErr := database.CreateRedisAccess(dbUser.ID, token)
	if utils.LogError(redisErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
	return
}

func refresh(c *gin.Context) {
	// Gets the refresh token from context and handle any possible errors
	refreshUUID, userID, getErr := middleware.GetRefreshTokenData(c)
	if utils.LogError(getErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Deletes the cached accesses in Redis
	redisErr := database.DeleteRedisAccess(refreshUUID)
	if utils.LogError(redisErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.RefreshError)
		return
	}

	// Creates new authentication tokens
	token, tokenErr := middleware.CreateToken(userID)
	if utils.LogError(tokenErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.RefreshError)
		return
	}

	// Sets the cached accesses in Redis
	redisErr = database.CreateRedisAccess(userID, token)
	if utils.LogError(redisErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.RefreshError)
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
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Finds a model in collection with the same inputted e-mail
	_, findErr := database.FindOne(database.UserCollection, bson.M{"email": user.Email})
	if utils.LogError(findErr) == nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseAlreadyExists)
		return
	}

	// Hashes and salts the user password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Declares new properties
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Inserts model into collection
	insertErr := database.InsertOne(database.UserCollection, user)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
	return
}

func logout(c *gin.Context) {
	// Gets the token metadata from request and handle any possible errors
	token, getErr := middleware.GetTokenData(c.Request)
	if utils.LogError(getErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	// Deletes the cached access from Redis and handle any possible errors
	delErr := database.DeleteRedisAccess(token.AccessUUID)
	if delErr != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.LogoutSuccess)
	return
}
