package database

import (
	"context"
	"errors"
	"fmt"
	"go_api_boilerplate/models"
	"go_api_boilerplate/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/redis/go-redis/v9"
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
	_, pingErr := client.Ping(context.Background()).Result()
	utils.LogFatalError(pingErr)

	log.Println("Redis client has been connected.")
}

// CreateRedisAccess expects an ID and a Token model
// to create the cached access
func CreateRedisAccess(ctx context.Context, id primitive.ObjectID, t *models.Token) error {
	// Gathers system times
	accessTime := time.Unix(t.AccessExpires, 0)
	refreshTime := time.Unix(t.RefreshExpires, 0)
	currentTime := time.Now()

	// Creates access token and handles any possible errors
	err := client.Set(ctx, t.AccessUUID, id.Hex(), accessTime.Sub(currentTime)).Err()
	if err != nil {
		return err
	}

	// Creates refresh token and handles any possible errors
	err = client.Set(ctx, t.RefreshUUID, id.Hex(), refreshTime.Sub(currentTime)).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetRedisAccess expects a RedisAccess model
// to return whether cached access is valid or not
func GetRedisAccess(ctx context.Context, access *models.RedisAccess) error {
	// Gathers the access from Redis and handles any possible errors
	_, err := client.Get(ctx, access.AccessUUID).Result()
	if err != nil {
		return err
	}

	return nil
}

// DeleteRedisAccess expects a RedisAccess UUID
// to remove the cached access
func DeleteRedisAccess(ctx context.Context, uuid string) error {
	// Deletes the cached access from Redis and handles any possible errors
	count, err := client.Del(ctx, uuid).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("redis: no keys in result")
	}

	return nil
}
