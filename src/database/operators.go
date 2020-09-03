package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeleteOne expects a collection and an ID in order to delete
// a document from the database
func DeleteOne(collection *mongo.Collection, id primitive.ObjectID) error {
	// Deletes a model from the database and handles any possible errors
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil || res.DeletedCount == 0 {
		return errors.New("mongo: no documents in result")
	}

	return nil
}

// FindAll expects a collection in order to find
// all documents into the database
func FindAll(collection *mongo.Collection) ([]bson.M, error) {
	// Creates a slice of documents
	var models []bson.M

	// Finds all models in the database and handles any possible errors
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	// Iterates over the cursor, i.e., list of documents
	for cursor.Next(context.Background()) {
		// Creates a document
		var model bson.M

		// Decodes the document and handles any possible errors
		err := cursor.Decode(&model)
		if err != nil {
			return nil, err
		}

		// Appends the model to the list
		models = append(models, model)
	}

	// Closes the iterator
	cursor.Close(context.Background())

	return models, nil
}

// FindAllWithAggregate expects a collection and an aggregation slice in order to find
// all documents into the database
func FindAllWithAggregate(collection *mongo.Collection, pipeline []bson.M) ([]bson.M, error) {
	// Creates a slice of documents
	var models []bson.M

	// Finds all models in the database and handles any possible errors
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	// Iterates over the cursor, i.e., list of documents
	for cursor.Next(context.Background()) {
		// Creates a document
		var model bson.M

		// Decodes the document and handles any possible errors
		err := cursor.Decode(&model)
		if err != nil {
			return nil, err
		}

		// Appends the model to the list
		models = append(models, model)
	}

	// Closes the iterator
	cursor.Close(context.Background())

	return models, nil
}

// FindOne expects a collection and a filter in order to find
// a single document into the database
func FindOne(collection *mongo.Collection, filter bson.M) (bson.M, error) {
	// Creates a document
	var model bson.M

	// Finds a model in the database and handles any possible errors
	// Note that it returns `nil` if model has been found
	err := collection.FindOne(context.Background(), filter).Decode(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// FindOneWithAggregate expects a collection and a pipeline in order to find
// a single document into the database
func FindOneWithAggregate(collection *mongo.Collection, pipeline []bson.M) (bson.M, error) {
	// Creates a document
	var model bson.M

	// Finds a model in the database and handles any possible errors
	// Note that it returns `nil` if model has been found
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	// Iterates over the cursor, i.e., list of documents
	for cursor.Next(context.Background()) {
		// Decodes the document and handles any possible errors
		err := cursor.Decode(&model)
		if err != nil {
			return nil, err
		}

	}

	// Closes the iterator
	cursor.Close(context.Background())

	return model, nil
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

// UpdateOne expects a collection, a filter and an update object in order to update
// a document from the database
func UpdateOne(collection *mongo.Collection, id primitive.ObjectID, update bson.M) error {
	// Updates a model from the database and handles any possible errors
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil || res.ModifiedCount == 0 {
		return errors.New("mongo: no documents in result")
	}

	return nil
}
