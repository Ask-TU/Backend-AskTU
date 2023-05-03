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

		questionId := c.Param("question_id")

		objId, err := primitive.ObjectIDFromHex(questionId)
		
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
			Owner_name:  answer.Owner_name,
			Question_id: questionId,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}

		result, err := answerCollection.InsertOne(ctx, newAnswer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result)

		var oldQuestion models.Question

		err = QeustionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&oldQuestion)
		
		Answers := append(oldQuestion.Answer, newAnswer.ID.Hex())

		update := bson.M{
			"content":    oldQuestion.Content,
			"owner":      oldQuestion.Class_id,
			"class_id":    oldQuestion.Class_id,
			"created_at":  oldQuestion.Created_at,
			"updated_at": time.Now(),
			"answer":     Answers,
		}

		question_result, err := QeustionCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			fmt.Println(question_result)
		}

		newNotifiactionForOwner := models.Notification{
			ID:           primitive.NewObjectID(),
			Content:	  "You are have created new answers",
			Owner:		answer.Owner,
			Class_id:   oldQuestion.Class_id,     
			Created_at:   time.Now(),
		}

		newNotifiactionForAnswers := models.Notification{
			ID:           primitive.NewObjectID(),
			Content:	  "You are have a new answers",
			Owner:		oldQuestion.Owner,
			Class_id:   oldQuestion.Class_id,     
			Created_at:   time.Now(),
		}
		result_notificationOn, err := NotificationCollection.InsertOne(ctx, newNotifiactionForOwner)
		result_notificationAn, err := NotificationCollection.InsertOne(ctx, newNotifiactionForAnswers)
		
		fmt.Println(result_notificationOn, result_notificationAn)

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": newAnswer}})

	}
}
func GetAllAnswers() gin.HandlerFunc {
	return func(c *gin.Context) {

		questionId := c.Param("question_id")

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
