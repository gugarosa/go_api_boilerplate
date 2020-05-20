package middleware

import (
	"os"
	"time"
	"vivere_api/models"
	"vivere_api/utils"

	"github.com/twinj/uuid"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateNewToken ...
func CreateNewToken(id primitive.ObjectID) (*models.Token, error) {
	t := &models.Token{}

	t.AccessExpires = time.Now().Add(time.Minute * 15).Unix()
	t.AccessUUID = uuid.NewV4().String()
	t.RefreshExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	t.RefreshUUID = uuid.NewV4().String()

	var err error

	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["access_uuid"] = t.AccessUUID
	claims["id"] = id
	claims["exp"] = t.AccessExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if utils.HandleErrors(err) != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = t.RefreshUUID
	rtClaims["id"] = id
	rtClaims["exp"] = t.RefreshExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if utils.HandleErrors(err) != nil {
		return nil, err
	}

	return t, nil
}
