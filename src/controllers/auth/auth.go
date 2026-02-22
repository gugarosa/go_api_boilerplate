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
	var inputUser, dbUser models.User

	checkErr := controllers.BindAndValidateRequest(c, &inputUser)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Use 401 for both wrong email and wrong password to prevent user enumeration
	user, findErr := database.FindOne(c.Request.Context(), database.UserCollection, bson.M{"email": inputUser.Email})
	if utils.LogError(findErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	decodeErr := controllers.EncodeStruct(user, &dbUser)
	if utils.LogError(decodeErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password))
	if utils.LogError(passwordErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.LoginError)
		return
	}

	token, tokenErr := middleware.CreateToken(dbUser.ID)
	if utils.LogError(tokenErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.LoginError)
		return
	}

	redisErr := database.CreateRedisAccess(c.Request.Context(), dbUser.ID, token)
	if utils.LogError(redisErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.LoginError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func refresh(c *gin.Context) {
	refreshUUID, userID, getErr := middleware.GetRefreshTokenData(c)
	if utils.LogError(getErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	redisErr := database.DeleteRedisAccess(c.Request.Context(), refreshUUID)
	if utils.LogError(redisErr) != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.RefreshError)
		return
	}

	token, tokenErr := middleware.CreateToken(userID)
	if utils.LogError(tokenErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RefreshError)
		return
	}

	redisErr = database.CreateRedisAccess(c.Request.Context(), userID, token)
	if utils.LogError(redisErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RefreshError)
		return
	}

	utils.DynamicResponse(c, http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func register(c *gin.Context) {
	var user models.User

	checkErr := controllers.BindAndValidateRequest(c, &user)
	if utils.LogError(checkErr) != nil {
		utils.ConstantResponse(c, http.StatusBadRequest, utils.RequestError)
		return
	}

	// Check for existing user with same email
	_, findErr := database.FindOne(c.Request.Context(), database.UserCollection, bson.M{"email": user.Email})
	if utils.LogError(findErr) == nil {
		utils.ConstantResponse(c, http.StatusConflict, utils.DatabaseAlreadyExists)
		return
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if utils.LogError(hashErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.RequestError)
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	insertErr := database.InsertOne(c.Request.Context(), database.UserCollection, user)
	if utils.LogError(insertErr) != nil {
		utils.ConstantResponse(c, http.StatusInternalServerError, utils.DatabaseInsertError)
		return
	}

	utils.ConstantResponse(c, http.StatusCreated, utils.DatabaseInsertSuccess)
}

func logout(c *gin.Context) {
	// Retrieve token data stored by AuthGuard middleware
	val, exists := c.Get("token_data")
	if !exists {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	token, ok := val.(*models.RedisAccess)
	if !ok {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	delErr := database.DeleteRedisAccess(c.Request.Context(), token.AccessUUID)
	if delErr != nil {
		utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
		return
	}

	utils.ConstantResponse(c, http.StatusOK, utils.LogoutSuccess)
}
