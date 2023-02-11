package controllers

import (
	"context"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"exmaple/Backendasktu/database"
	"exmaple/Backendasktu/responses"

	"exmaple/Backendasktu/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var answerCollection *mongo.Collection = database.OpenCollection(database.Client, "answers")

func CreateAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {

		qestionId := c.Param("questionId")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var answer models.Answer

		if err := c.BindJSON(&answer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newAnswer := models.Answer{
			ID:          primitive.NewObjectID(),
			Content:     answer.Content,
			Owner:       answer.Owner,
			Question_id: qestionId,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}

		result, err := answerCollection.InsertOne(ctx, newAnswer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result)

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": newAnswer}})

	}
}
func GetAllAnswers() gin.HandlerFunc {
	return func(c *gin.Context) {

		questionId := c.Param("questionId")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		answers, err := findAnswersByClassId(ctx, questionId)
		if err != nil {
			c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Fail to Find Data", Result: map[string]interface{}{"data": err.Error()}})
			return
		}
		fmt.Println(answers)
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": answers}})
	}
}

func findAnswersByClassId(ctx context.Context, qestionId string) ([]interface{}, error) {
	var answers []interface{}

	results, err := answerCollection.Find(ctx, bson.M{"question_id": qestionId})
	if err != nil {
		return nil, err
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleAnswer models.Answer
		if err = results.Decode(&singleAnswer); err != nil {
			return nil, err
		}
		log.Println(singleAnswer)
		answers = append(answers, singleAnswer)
	}

	return answers, nil
}

func GetAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		answerID := c.Param("answerID")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var answer1 models.Answer
		AnsId, _ := primitive.ObjectIDFromHex(answerID)

		err := answerCollection.FindOne(ctx, bson.M{"_id": AnsId}).Decode(&answer1)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, answer1)
		fmt.Println(answer1)
	}
}
