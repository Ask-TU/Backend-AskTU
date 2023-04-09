package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID         primitive.ObjectID `bson:"_id"`
	Content    string             `json:"Content"`
	Owner      string             `json:"owner"`
	Class_id   string             `json:"class_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Answer     []string           `json:"answer"`
}
