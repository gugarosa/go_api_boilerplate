package database

import (
	"context"
	"go_api_boilerplate/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// InitMongo expects an url connection string and a database name
// in order to start a MongoDB service
func InitMongo(url string, database string) {
	// Creates a context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connects to MongoDB and handles any possible errors
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	utils.LogFatalError(err)

	// Pings client to check its connection and handles possible fatal error
	pingErr := client.Ping(ctx, readpref.Primary())
	utils.LogFatalError(pingErr)

	// If no error has occurred, just log that the client has been connected
	log.Println(utils.DatabaseClientConnected)

	// Retrieves the database
	dbObj := client.Database(database)

	// Adds desired collections
	SetCollections(dbObj)
}
