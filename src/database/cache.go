package database

import (
	"errors"
	"fmt"
	"go_api_boilerplate/models"
	"go_api_boilerplate/utils"
	"time"

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
	utils.LogFatalError(pingErr)

}

// CreateRedisAccess expects an ID and a Token model
// to create the cached access
func CreateRedisAccess(id primitive.ObjectID, t *models.Token) error {
	// Gathers system times
	accessTime := time.Unix(t.AccessExpires, 0)
	refreshTime := time.Unix(t.RefreshExpires, 0)
	currentTime := time.Now()

	// Creates access token and handles any possible errors
	err := client.Set(t.AccessUUID, id.Hex(), accessTime.Sub(currentTime)).Err()
	if err != nil {
		return err
	}

	// Creates refresh token and handles any possible errors
	err = client.Set(t.RefreshUUID, id.Hex(), refreshTime.Sub(currentTime)).Err()
	if err != nil {
		return err
	}

	return nil
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
	count, err := client.Del(uuid).Result()
	if err != nil || count == 0 {
		return errors.New("redis: no keys in result")
	}

	return nil
}
