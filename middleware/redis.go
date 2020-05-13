package middleware

import (
	"os"
	"time"
	"vivere_api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-redis/redis/v7"
)

// Client
var client *redis.Client

// InitializeRedis ...
func InitializeRedis() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
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
