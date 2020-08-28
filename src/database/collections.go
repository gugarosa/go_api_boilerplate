package database

import "go.mongodb.org/mongo-driver/mongo"

// ProductCollection global variable
var ProductCollection *mongo.Collection

// UserCollection global variable
var UserCollection *mongo.Collection

// SetCollections expects a MongoDB database as parameter and sets in the scope
// variables to the desired collections
func SetCollections(c *mongo.Database) {
	ProductCollection = c.Collection("products")
	UserCollection = c.Collection("users")
}
