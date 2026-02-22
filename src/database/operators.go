package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeleteOne expects a context, collection and an ID in order to delete
// a document from the database
func DeleteOne(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) error {
	// Deletes a model from the database and handles any possible errors
	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("mongo: no documents in result")
	}

	return nil
}

// FindAll expects a context and collection in order to find
// all documents into the database
func FindAll(ctx context.Context, collection *mongo.Collection) ([]bson.M, error) {
	models := make([]bson.M, 0)

	// Finds all models in the database and handles any possible errors
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decodes all documents from the cursor
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// FindAllWithAggregate expects a context, collection and an aggregation slice in order to find
// all documents into the database
func FindAllWithAggregate(ctx context.Context, collection *mongo.Collection, pipeline []bson.M) ([]bson.M, error) {
	models := make([]bson.M, 0)

	// Finds all models in the database and handles any possible errors
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decodes all documents from the cursor
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// FindOne expects a context, collection and a filter in order to find
// a single document into the database
func FindOne(ctx context.Context, collection *mongo.Collection, filter bson.M) (bson.M, error) {
	var model bson.M

	// Finds a model in the database and handles any possible errors
	err := collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// FindOneWithAggregate expects a context, collection and a pipeline in order to find
// a single document into the database
func FindOneWithAggregate(ctx context.Context, collection *mongo.Collection, pipeline []bson.M) (bson.M, error) {
	var models []bson.M

	// Finds a model in the database and handles any possible errors
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decodes all documents from the cursor
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	if len(models) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return models[0], nil
}

// InsertOne expects a context, collection and a model in order to insert
// a new document into the database
func InsertOne(ctx context.Context, collection *mongo.Collection, model interface{}) error {
	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOne expects a context, collection, an ID and an update object in order to update
// a document from the database
func UpdateOne(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, update bson.M) error {
	res, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("mongo: no documents in result")
	}

	return nil
}
