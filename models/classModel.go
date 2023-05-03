package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Classrooms struct {
	ID           primitive.ObjectID `bson:"_id"`
	Subject_name string             `json:"subject_name"`
	Class_owner  string             `json:"class_owner"`
	Tag 		 string             `json:"tag"`
	Section      string             `json:"section"`
	Class_id     string             `json:"class_id"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	Questions    []string           `json:"questions"`
	Members      []string           `json:"members"`
}
