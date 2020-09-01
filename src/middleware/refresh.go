package middleware

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetRefreshToken expects an incoming context
// to get the refresh token from body
func GetRefreshToken(c *gin.Context) string {
	// Creates a request variable
	request := map[string]string{}

	// Binds the model and handle any possible errors
	err := c.ShouldBindJSON(&request)
	if err != nil {
		return ""
	}

	// Gathers the refresh token from request
	refreshToken := request["refresh_token"]

	return refreshToken
}

// VerifyRefreshToken expects an incoming context
// to verify whether refresh token signature is valid
func VerifyRefreshToken(c *gin.Context) (*jwt.Token, error) {
	// Gathers the refresh token
	refreshToken := GetRefreshToken(c)

	// Parses the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Ensures that it conforms to `SigningMethodHMAC` and handles any possible errors
		_, valid := token.Method.(*jwt.SigningMethodHMAC)
		if !valid {
			return nil, fmt.Errorf("Unexpected signing: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	// Handles any possible errors
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetRefreshTokenData expects an incoming context
// to get refresh token meta-data
func GetRefreshTokenData(c *gin.Context) (string, primitive.ObjectID, error) {
	// Gathers an already-verified token and handle any possible errors
	token, err := VerifyRefreshToken(c)
	if err != nil {
		return "", primitive.NilObjectID, err
	}

	// Gathers the refresh token metadata and its validation
	claims, valid := token.Claims.(jwt.MapClaims)

	// Checks whether refresh token is valid and handles any possible errors
	if valid && token.Valid {
		refreshUUID, valid := claims["refresh_uuid"].(string)
		userID, err := primitive.ObjectIDFromHex(claims["id"].(string))

		if !valid || err != nil {
			return "", primitive.NilObjectID, err
		}

		return refreshUUID, userID, nil
	}

	return "", primitive.NilObjectID, err

}
