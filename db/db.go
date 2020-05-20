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

// InitDatabase expects an url connection string and a database name
// in order to start a MongoDB service
func InitDatabase(url string, database string) {
	// Creating a new MongoDB client, setting context and connecting
	client, _ := mongo.NewClient(options.Client().ApplyURI(url))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = client.Connect(ctx)

	// Closes the connection at the ending
	defer cancel()

	// Pinging client to check its connection and handling possible eror
	pingErr := client.Ping(context.Background(), readpref.Primary())
	utils.HandleFatalError(pingErr)

	// If no error has occured, just log that the client has been connected
	log.Println(utils.DatabaseClientConnected)

	// Retrieving the database
	dbObj := client.Database(database)

	// Adding desired collections
	SetCollections(dbObj)

	return
}
