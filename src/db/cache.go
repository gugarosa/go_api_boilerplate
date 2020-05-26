package db

import (
	"errors"
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
func InitRedis(host string, port string, password string) {
	// Creates a Redis client
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	// Pings client to check its connection and handles possible fatal error
	_, pingErr := client.Ping().Result()
	utils.HandleFatalError(pingErr)

}

// CreateRedisAccess expects an ID and a Token model
// to create the cached access
func CreateRedisAccess(id primitive.ObjectID, t *models.Token) (error, error) {
	// Gathers system times
	accessTime := time.Unix(t.AccessExpires, 0)
	refreshTime := time.Unix(t.RefreshExpires, 0)
	currentTime := time.Now()

	// Creates access and refresh tokens
	err := client.Set(t.AccessUUID, id.Hex(), accessTime.Sub(currentTime)).Err()
	err2 := client.Set(t.RefreshUUID, id.Hex(), refreshTime.Sub(currentTime)).Err()

	return err, err2
}

// GetRedisAccess expects a RedisAccess model
// to return whether cached access is valid or not
func GetRedisAccess(access *models.RedisAccess) error {
	// Gathers the access from Redis and handles any possible errors
	_, err := client.Get(access.AccessUUID).Result()
	if err != nil {
		return err
	}

	return nil
}

// DeleteRedisAccess expects a RedisAccess UUID
// to remove the cached access
func DeleteRedisAccess(uuid string) error {
	// Deletes the cached access from Redis and handles any possible errors
	amountKey, err := client.Del(uuid).Result()
	if err != nil || amountKey == 0 {
		return errors.New("key could not be deleted")
	}

	return nil
}
