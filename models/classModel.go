package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllClass struct {
	ID           primitive.ObjectID `bson:"_id"`
	Main_id      string             `json:"main_id"`
	Subject_name string             `json:"subject_name"`
	Class_owner  string             `json:"class_owner"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	Question     []Question         `json:"question"`
	Members      []Member           `json:"member"`
}

type Question struct {
	ID         primitive.ObjectID `bson:"_id"`
	Content    string             `json:"question"`
	Owner      string             `json:"owner"`
	Class_id   string             `json:"class_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Answer     []Answer           `json:"answer"`
}

type Answer struct {
	ID          primitive.ObjectID `bson:"_id"`
	Content     string             `json:"answer"`
	Owner       string             `json:"owner"`
	Question_id string             `json:"question_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}

type Member struct {
	ID        primitive.ObjectID `bson:"_id"`
	ClassRoom string             `json:"classroom"`
}
