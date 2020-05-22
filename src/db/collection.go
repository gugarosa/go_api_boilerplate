package db

import "go.mongodb.org/mongo-driver/mongo"

// UserCollection global variable
var UserCollection *mongo.Collection

// SetCollections expects a MongoDB database as parameter and sets in the scope
// variables to the desired collections
func SetCollections(c *mongo.Database) {
	UserCollection = c.Collection("users")
}
