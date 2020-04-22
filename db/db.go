package db

import (
	"context"
	"log"
	"time"
	"vivere_api/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// InitializeDatabase starts a MongoDB service on desired URL
func InitializeDatabase(url string, database string) {
	// Creating a new MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	// Setting up a context, required by Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	// Closes the connection at the ending
	defer cancel()

	// Pinging client to check its connection
	err = client.Ping(context.Background(), readpref.Primary())

	// Handling if there is a fatal error
	utils.HandleFatalError(err)

	// If no error has occured, just log that the client has been connected
	log.Println("MongoDB client connected.")

	// Retrieving the database
	_db := client.Database(database)

	// Adding desired collections
	SetCollections(_db)

	return
}
