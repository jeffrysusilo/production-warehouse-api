package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductionLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductionID primitive.ObjectID `bson:"production_id" json:"production_id"`
	Action       string             `bson:"action" json:"action"` // created, completed, canceled
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
	Note         string             `bson:"note,omitempty" json:"note,omitempty"`
}
