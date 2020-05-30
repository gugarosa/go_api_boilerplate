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
	// Creates a new MongoDB client, sets context and connects
	client, _ := mongo.NewClient(options.Client().ApplyURI(url))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = client.Connect(ctx)

	// Closes the connection at the ending
	defer cancel()

	// Pings client to check its connection and handles possible fatal error
	pingErr := client.Ping(context.Background(), readpref.Primary())
	utils.LogFatalError(pingErr)

	// If no error has occured, just log that the client has been connected
	log.Println(utils.DatabaseClientConnected)

	// Retrieves the database
	dbObj := client.Database(database)

	// Adds desired collections
	SetCollections(dbObj)

	return
}
