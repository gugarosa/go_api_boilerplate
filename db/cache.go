package db

import (
	"fmt"
	"time"
	"vivere_api/models"
	"vivere_api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-redis/redis/v7"
)

// Client
var client *redis.Client

// InitRedis expects a port and a password
// to initialize a Redis service
func InitRedis(port string, password string) {
	// Creating a Redis client
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("cache:%s", port),
		Password: password,
		DB:       0,
	})

	// Pinging client to check its connection and handling possible eror
	_, pingErr := client.Ping().Result()
	utils.HandleFatalError(pingErr)

}

// CreateAuth ...
func CreateAuth(id primitive.ObjectID, t *models.Token) error {
	at := time.Unix(t.AccessExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(t.RefreshExpires, 0)
	now := time.Now()

	errAccess := client.Set(t.AccessUUID, id.Hex(), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(t.RefreshUUID, id.Hex(), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
