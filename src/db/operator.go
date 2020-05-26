package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOne expects a collection, a model and a filter in order to find
// a single document into the database
func FindOne(collection *mongo.Collection, filter bson.M, model interface{}) error {
	// Tries to find a model in the database
	// Note that it returns `nil` if model has been found
	err := collection.FindOne(context.Background(), filter).Decode(model)

	return err
}

// InsertOne expects a collection and a model in order to insert
// a new document into the database
func InsertOne(collection *mongo.Collection, model interface{}) error {
	// Tries to insert the model in the database
	_, err := collection.InsertOne(context.Background(), model)

	return err
}
