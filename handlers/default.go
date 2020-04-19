package handlers

import "go.mongodb.org/mongo-driver/mongo"

// Global variables for the DefaultHandler
var userCollection *mongo.Collection

// SetCollections expects a MongoDB database as parameter and sets in the scope
// variables to the desired collections
func SetCollections(c *mongo.Database) {
	userCollection = c.Collection("users")
}
