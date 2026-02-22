package middleware

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetRefreshToken expects an incoming context
// to get the refresh token from body
func GetRefreshToken(c *gin.Context) string {
	request := map[string]string{}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		return ""
	}

	return request["refresh_token"]
}

// VerifyRefreshToken expects an incoming context
// to verify whether refresh token signature is valid
func VerifyRefreshToken(c *gin.Context) (*jwt.Token, error) {
	refreshToken := GetRefreshToken(c)

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		_, valid := token.Method.(*jwt.SigningMethodHMAC)
		if !valid {
			return nil, fmt.Errorf("unexpected signing: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetRefreshTokenData expects an incoming context
// to get refresh token meta-data
func GetRefreshTokenData(c *gin.Context) (string, primitive.ObjectID, error) {
	token, err := VerifyRefreshToken(c)
	if err != nil {
		return "", primitive.NilObjectID, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", primitive.NilObjectID, errors.New("invalid refresh token claims")
	}

	refreshUUID, ok := claims["refresh_uuid"].(string)
	if !ok {
		return "", primitive.NilObjectID, errors.New("invalid refresh_uuid claim")
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		return "", primitive.NilObjectID, errors.New("invalid id claim")
	}

	userID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return "", primitive.NilObjectID, err
	}

	return refreshUUID, userID, nil
}
