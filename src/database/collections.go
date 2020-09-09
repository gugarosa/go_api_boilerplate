package database

import "go.mongodb.org/mongo-driver/mongo"

// CategoryCollection global variable
var CategoryCollection *mongo.Collection

// ProductCollection global variable
var ProductCollection *mongo.Collection

// QuestionCollection global variable
var QuestionCollection *mongo.Collection

// SurveyCollection global variable
var SurveyCollection *mongo.Collection

// TagCollection global variable
var TagCollection *mongo.Collection

// UserCollection global variable
var UserCollection *mongo.Collection

// SetCollections expects a MongoDB database as parameter and sets in the scope
// variables to the desired collections
func SetCollections(c *mongo.Database) {
	CategoryCollection = c.Collection("categories")
	ProductCollection = c.Collection("products")
	QuestionCollection = c.Collection("questions")
	SurveyCollection = c.Collection("surveys")
	TagCollection = c.Collection("tags")
	UserCollection = c.Collection("users")
}
