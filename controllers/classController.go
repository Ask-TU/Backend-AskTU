package controllers

import (
	"context"
	"fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"exmaple/Backendasktu/database"

	"exmaple/Backendasktu/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var classroomCollection *mongo.Collection = database.OpenCollection(database.Client, "classrooms")

var class = []models.AllClass{
	{Main_id: "1",
		Subject_name: "SF230",
		Class_owner:  "UserID_1",
		Question: []models.Question{
			{Class_id: "1", Content: "What is GoLang", Owner: "UserID_1", Answer: []models.Answer{
				{Question_id: "1", Content: "Golang is 1", Owner: "UserID_2"},
				{Question_id: "1", Content: "Golang is 2", Owner: "UserID_3"},
			}},
			{Class_id: "1", Content: "What is Java", Owner: "UserID_2", Answer: []models.Answer{
				{Question_id: "2", Content: "Java is 1", Owner: "UserID_3"},
				{Question_id: "2", Content: "Java is 2", Owner: "UserID_4"},
			}},
		}},

	{Main_id: "2",
		Subject_name: "SF555",
		Class_owner:  "UserID_1",
		Question: []models.Question{
			{Class_id: "2", Content: "What is HTML", Owner: "UserID_1", Answer: []models.Answer{
				{Question_id: "3", Content: "HTML is 1", Owner: "UserID_2"},
				{Question_id: "3", Content: "HTML is 2", Owner: "UserID_3"},
			}},
			{Class_id: "2", Content: "What is Lue", Owner: "UserID_2", Answer: []models.Answer{
				{Question_id: "4", Content: "Lue is 1", Owner: "UserID_3"},
				{Question_id: "4", Content: "Lue is 2", Owner: "UserID_4"},
			}},
		}},
}

func GetClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, class)
	}
}

func GetQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, class[0].Question[0:])
	}
}

func CreateClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		newUser := models.AllClass{
			ID:           primitive.NewObjectID(),
			Subject_name: "eiei",
			Class_owner:  "user.Name",
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
			Question:     []models.Question{},
			Members:      []models.Member{},
		}

		result, err := classroomCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result)
		c.JSON(http.StatusCreated, "success")

	}
}

func DeleteClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		for i, a := range class {
			if a.Main_id == id {
				class = append(class[:i], class[i+1:]...)
				break
			}
		}

		c.Status(http.StatusNoContent)
	}
}
