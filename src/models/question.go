package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question model
type Question struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Description string               `bson:"description" json:"description" validate:"required"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	Active      bool                 `bson:"active,omitempty"`
	CreatedAt   time.Time            `bson:"created_at,omitempty"`
	UpdatedAt   time.Time            `bson:"updated_at,omitempty"`
}
