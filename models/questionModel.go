package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID         primitive.ObjectID `bson:"_id"`
	Content    string             `json:"Content"`
	Owner      string             `json:"owner"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Answer     []Answer           `json:"answer"`
}

type Answer struct {
	ID          primitive.ObjectID `bson:"_id"`
	Content     string             `json:"content"`
	Owner       string             `json:"owner"`
	Question_id string             `json:"question_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
