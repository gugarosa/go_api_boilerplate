package db

import (
	"fmt"
	"time"
	"vivere_api/models"
	"vivere_api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-redis/redis/v7"
)

// Redis client variable
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

	// Pinging client to check its connection and handling its possible fatal error
	_, pingErr := client.Ping().Result()
	utils.HandleFatalError(pingErr)

}

// SetTokens expects an ID and a Token model
// to set the access and refresh tokens
func SetTokens(id primitive.ObjectID, t *models.Token) bool {
	// Gathering times
	accessTime := time.Unix(t.AccessExpires, 0)
	refreshTime := time.Unix(t.RefreshExpires, 0)
	currentTime := time.Now()

	// Trying to set access and refresh token
	accessErr := client.Set(t.AccessUUID, id.Hex(), accessTime.Sub(currentTime)).Err()
	refreshErr := client.Set(t.RefreshUUID, id.Hex(), refreshTime.Sub(currentTime)).Err()

	// Handling and returning any possible errors
	return utils.HandleError(accessErr, refreshErr)
}
