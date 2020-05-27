package db

import "go.mongodb.org/mongo-driver/mongo"

// ItemCollection global variable
var ItemCollection *mongo.Collection

// UserCollection global variable
var UserCollection *mongo.Collection

// SetCollections expects a MongoDB database as parameter and sets in the scope
// variables to the desired collections
func SetCollections(c *mongo.Database) {
	ItemCollection = c.Collection("items")
	UserCollection = c.Collection("users")
}
