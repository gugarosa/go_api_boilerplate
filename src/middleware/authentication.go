package middleware

import (
	"errors"
	"fmt"
	"go_api_boilerplate/models"
	"go_api_boilerplate/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateToken expects an user identifier
// to create and sign the JWT tokens
func CreateToken(id primitive.ObjectID) (*models.Token, error) {
	token := &models.Token{}

	// Defines the access token meta-information
	token.AccessExpires = time.Now().Add(time.Minute * 15).Unix()
	token.AccessUUID = uuid.New().String()

	// Defines the refresh token meta-information
	token.RefreshExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	token.RefreshUUID = uuid.New().String()

	// Creating the access token structure and signing it with JWT
	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["access_uuid"] = token.AccessUUID
	accessClaims["id"] = id.Hex()
	accessClaims["exp"] = token.AccessExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	var err error

	token.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating the refresh token structure and signing it with JWT
	refreshClaims := jwt.MapClaims{}
	refreshClaims["refresh_uuid"] = token.RefreshUUID
	refreshClaims["id"] = id.Hex()
	refreshClaims["exp"] = token.RefreshExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	token.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetToken expects an incoming request
// to get the access token
func GetToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	splitToken := strings.Split(token, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}

	return ""
}

// VerifyToken expects an incoming request
// to verify whether token signature is valid
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenStr := GetToken(r)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, valid := token.Method.(*jwt.SigningMethodHMAC)
		if !valid {
			return nil, fmt.Errorf("unexpected signing: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetTokenData expects an incoming request
// to gather the token metadata
func GetTokenData(r *http.Request) (*models.RedisAccess, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, errors.New("invalid access_uuid claim")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return nil, errors.New("invalid id claim")
	}

	return &models.RedisAccess{
		AccessUUID: accessUUID,
		UserID:     userID,
	}, nil
}

// AuthGuard provides an authentication guard
// to check whether token is valid or not
func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, err := GetTokenData(c.Request)
		if err != nil {
			utils.ConstantResponse(c, http.StatusUnauthorized, utils.AuthError)
			c.Abort()
			return
		}

		// Store parsed token data so handlers don't re-parse the JWT
		c.Set("token_data", tokenData)
		c.Next()
	}
}
