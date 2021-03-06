package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tag model
type Tag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}
