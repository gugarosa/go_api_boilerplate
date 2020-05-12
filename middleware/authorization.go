package middleware

import (
	"os"
	"time"
	"vivere_api/utils"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateNewToken ...
func CreateNewToken(id primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if !utils.HandleError(err) {
		return "", err
	}

	return token, nil
}
