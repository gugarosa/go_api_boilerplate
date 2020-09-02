package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Survey model
type Survey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Active    bool               `bson:"active,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}
