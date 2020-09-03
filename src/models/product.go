package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category model
type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

// Tag model
type Tag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

// Product model
type Product struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name" json:"name" validate:"required"`
	Brand       string               `bson:"brand" json:"brand" validate:"required"`
	Categories  []primitive.ObjectID `bson:"categories" json:"categories" validate:"required"`
	Summary     string               `bson:"summary" json:"summary"`
	Description string               `bson:"description" json:"description"`
	Image       string               `bson:"image" json:"image"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	Active      bool                 `bson:"active,omitempty"`
	CreatedAt   time.Time            `bson:"created_at,omitempty"`
	UpdatedAt   time.Time            `bson:"updated_at,omitempty"`
}
