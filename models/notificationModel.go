package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID          primitive.ObjectID `bson:"_id"`
	Content     string             `json:"content"`
	Owner       string             `json:"owner"`
	Class_id     string             `json:"class_id"`
	Created_at  time.Time          `json:"created_at"`
}

type Notifications struct {
	ID          primitive.ObjectID `bson:"_id"`
	Receive    string             `json:"receive"`
	Sender 	    string             `json:"sender"`
	Owner       string             `json:"owner"`
	Class_id    string             `json:"class_id"`
	Created_at  time.Time          `json:"created_at"`
}