package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID         primitive.ObjectID `bson:"_id"`
	Content    string             `json:"Content"`
	Owner      string             `json:"owner"`
	Owner_name   string            `json:"owner_name"`
	Tag  		string 			  `json:"tag"`
	Class_id   string             `json:"class_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Answer     []string           `json:"answer"`
}

type Answer struct {
	ID          primitive.ObjectID `bson:"_id"`
	Content     string             `json:"content"`
	Owner       string             `json:"owner"`
	Owner_name   string            `json:"owner_name"`
	Question_id string             `json:"question_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
