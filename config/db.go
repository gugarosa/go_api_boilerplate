package config

import (
	"context"
	"log"
	"time"
	"vivere_api/handlers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConfigureDatabase starts a MongoDB service on desired URL
func ConfigureDatabase(url string, database string) {
	// Creating a new MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	// Setting up a context, required by Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	// Closes the connection at the ending
	defer cancel()

	// Pinging client to check its connection
	err = client.Ping(context.Background(), readpref.Primary())

	// Checking if there is an error
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("MongoDB client connected.")
	}

	// Retrieving the database
	_db := client.Database(database)

	// Adding desired collections
	handlers.SetCollections(_db)

	return

}
