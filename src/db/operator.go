package db

import (
	"context"
	"vivere_api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOne expects a collection, a model and a filter in order to find
// a single document into the database
func FindOne(collection *mongo.Collection, filter bson.M, model interface{}) error {
	// Tries to find a model in the database
	// Note that it returns `nil` if model has been found
	findErr := collection.FindOne(context.Background(), filter).Decode(model)

	return utils.HandleError(findErr)
}

// InsertOne expects a collection and a model in order to insert
// a new document into the database
func InsertOne(collection *mongo.Collection, model interface{}) error {
	// Tries to insert the model in the database
	_, insertErr := collection.InsertOne(context.Background(), model)

	return utils.HandleError(insertErr)
}
