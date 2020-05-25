package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"vivere_api/models"
	"vivere_api/utils"

	"github.com/twinj/uuid"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateToken expects an user identifier
// to create and sign the JWT tokens
func CreateToken(id primitive.ObjectID) (*models.Token, error) {
	// Creates a new Token-based struct
	t := &models.Token{}

	// Defines the access token meta-information
	t.AccessExpires = time.Now().Add(time.Minute * 15).Unix()
	t.AccessUUID = uuid.NewV4().String()

	// Defines the refresh token meta-information
	t.RefreshExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	t.RefreshUUID = uuid.NewV4().String()

	// Creating the access token structure and signing it with JWT
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["access_uuid"] = t.AccessUUID
	claims["id"] = id.Hex()
	claims["exp"] = t.AccessExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Defines a local-scoped variable for handling the error
	var signErr error

	// Creates the access token and handle any possible errors
	t.AccessToken, signErr = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if utils.HandleError(signErr) != nil {
		return nil, signErr
	}

	// Creating the refresh token structure and signing it with JWT
	claims = jwt.MapClaims{}
	claims["refresh_uuid"] = t.RefreshUUID
	claims["id"] = id
	claims["exp"] = t.RefreshExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Creates the refresh token and handles any possible errors
	t.RefreshToken, signErr = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if utils.HandleError(signErr) != nil {
		return nil, signErr
	}

	return t, nil
}

// GetRequestToken ...
func GetRequestToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyToken ...
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := GetRequestToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ValidateToken ...
func ValidateToken(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// GetTokenData ...
func GetTokenData(r *http.Request) (*models.Access, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		// userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["id"]), 10, 64)
		userID := claims["id"].(string)
		return &models.Access{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}
