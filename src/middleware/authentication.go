package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"vivere_api/models"
	"vivere_api/utils"

	"github.com/gin-gonic/gin"

	"github.com/twinj/uuid"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateToken expects an user identifier
// to create and sign the JWT tokens
func CreateToken(id primitive.ObjectID) (*models.Token, error) {
	// Creates a new Token-based struct
	token := &models.Token{}

	// Defines the access token meta-information
	token.AccessExpires = time.Now().Add(time.Minute * 15).Unix()
	token.AccessUUID = uuid.NewV4().String()

	// Defines the refresh token meta-information
	token.RefreshExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	token.RefreshUUID = uuid.NewV4().String()

	// Creating the access token structure and signing it with JWT
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["access_uuid"] = token.AccessUUID
	claims["id"] = id.Hex()
	claims["exp"] = token.AccessExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Defines a local-scoped variable for handling the error
	var err error

	// Creates the access token and handle any possible errors
	token.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating the refresh token structure and signing it with JWT
	claims = jwt.MapClaims{}
	claims["refresh_uuid"] = token.RefreshUUID
	claims["id"] = id
	claims["exp"] = token.RefreshExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Creates the refresh token and handles any possible errors
	token.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetToken expects an incoming request
// to get the access token
func GetToken(r *http.Request) string {
	// Gathers the token from the `Authorization` header
	token := r.Header.Get("Authorization")

	// Splits the token into two parts
	// as it is composed of `Bearer <token>`
	splitToken := strings.Split(token, " ")

	// Checks if it could be splitted into two parts
	if len(splitToken) == 2 {
		return splitToken[1]
	}

	return ""
}

// VerifyToken expects an incoming request
// to verify whether token signature is valid
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	// Gathers the token
	tokenStr := GetToken(r)

	// Parses the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensures that it conforms to `SigningMethodHMAC` and handles any possible errors
		_, valid := token.Method.(*jwt.SigningMethodHMAC)
		if !valid {
			return nil, fmt.Errorf("Unexpected signing: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	// Handles any possible errors
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetTokenData expects an incoming request
// to gather the token metadata
func GetTokenData(r *http.Request) (*models.RedisAccess, error) {
	// Gathers an already-verified token and handle any possible errors
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}

	// Gathers the token metadata and its validation
	claims, valid := token.Claims.(jwt.MapClaims)

	// Checks whether token is valid and handle any possible errors
	if valid && token.Valid {
		accessUUID, valid := claims["access_uuid"].(string)
		if !valid {
			return nil, err
		}

		userID := claims["id"].(string)

		return &models.RedisAccess{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

// ValidateToken expects an incoming request
// to verify whether token is valid or not
func ValidateToken(r *http.Request) error {
	// Gathers an already-verified token and handle any possible errors
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}

	// Checks whether token valid and handle any possible errors
	_, valid := token.Claims.(jwt.Claims)
	if !valid && !token.Valid {
		return err
	}

	return nil
}

// AuthGuard provides an authentication guard
// to check whether token is valid or not
func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Checks if token is valid and handle any possible errors
		err := ValidateToken(c.Request)
		if err != nil {
			utils.StaticResponse(c, http.StatusUnauthorized, utils.AuthError)
			c.Abort()
			return
		}

		// If token is valid, proceed with the request
		c.Next()
	}
}
