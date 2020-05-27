package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindAll expects a collection and a model in order to find
// all documents into the database
func FindAll(collection *mongo.Collection) ([]bson.M, error) {
	var models []bson.M
	// Finds all models in the database and handles any possible errors
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var model bson.M
		err := cursor.Decode(&model)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	cursor.Close(context.Background())

	return models, nil
}

// FindOne expects a collection, a model and a filter in order to find
// a single document into the database
func FindOne(collection *mongo.Collection, filter bson.M, model interface{}) error {
	// Finds a model in the database and handles any possible errors
	// Note that it returns `nil` if model has been found
	err := collection.FindOne(context.Background(), filter).Decode(model)
	if err != nil {
		return err
	}

	return nil
}

// InsertOne expects a collection and a model in order to insert
// a new document into the database
func InsertOne(collection *mongo.Collection, model interface{}) error {
	// Inserts a model in the database and handles any possible errors
	_, err := collection.InsertOne(context.Background(), model)
	if err != nil {
		return err
	}

	return nil
}
