package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Warehouse   string             `bson:"warehouse" json:"warehouse"`
	Description string             `bson:"description" json:"description"`
}
