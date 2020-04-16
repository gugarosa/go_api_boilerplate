package db

import "go.mongodb.org/mongo-driver/mongo"

// Global variable for the user collection
var collection *mongo.Collection

// UserCollection expects a MongoDB database as parameter and sets in the scope
// a variable to the `users` collection
func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}
