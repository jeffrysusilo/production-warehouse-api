package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaterialUsed struct {
	ItemID        primitive.ObjectID `bson:"item_id" json:"item_id"`
	QuantityUsed  int                `bson:"quantity_used" json:"quantity_used"`
}

type Production struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductName      string             `bson:"product_name" json:"product_name"`
	Materials        []MaterialUsed     `bson:"materials" json:"materials"`
	QuantityProduced int                `bson:"quantity_produced" json:"quantity_produced"`
	ProductionDate   time.Time          `bson:"production_date" json:"production_date"`
}
