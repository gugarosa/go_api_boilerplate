package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product model
type Product struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name" json:"name" validate:"required"`
	Brand       string               `bson:"brand" json:"brand" validate:"required"`
	Category    string               `bson:"category" json:"category" validate:"required"`
	Summary     string               `bson:"summary" json:"summary"`
	Description string               `bson:"description" json:"description"`
	Image       string               `bson:"image" json:"image"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	Active      bool                 `bson:"active,omitempty"`
	CreatedAt   time.Time            `bson:"created_at,omitempty"`
	UpdatedAt   time.Time            `bson:"updated_at,omitempty"`
}
