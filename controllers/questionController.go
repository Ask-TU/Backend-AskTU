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

var QeustionCollection *mongo.Collection = database.OpenCollection(database.Client, "questions")

func CreateQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		classroom_id := c.Param("classroom_id")
		objId, err := primitive.ObjectIDFromHex(classroom_id)

		//var oldClass models.Classrooms
		var question models.Question

		defer cancel()
		if err := c.BindJSON(&question); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&question)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		newQuestion := models.Question{
			ID:         primitive.NewObjectID(),
			Content:    question.Content,
			Owner:      question.Owner,
			Owner_name: question.Owner_name,
			Class_id:   classroom_id,
			Created_at: time.Now(),
			Updated_at: time.Now(),
			Answer:     question.Answer,
		}
		
		result, err := QeustionCollection.InsertOne(ctx, newQuestion)

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result)


		var oldClass models.Classrooms
		err = ClassroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&oldClass)
		
		classroomQuestion := append(oldClass.Questions, newQuestion.ID.Hex())

		update := bson.M{
			"subject_name": oldClass.Subject_name,
			"class_owner":  oldClass.Class_owner,
			"created_at":   oldClass.Created_at,
			"updated_at":   time.Now(),
			"questions":    classroomQuestion,
			"members":      oldClass.Members,
		}

		classroom_result, err := ClassroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		
		newNotifiaction := models.Notification{
			ID:           primitive.NewObjectID(),
			Content:	  "You have created a new Question",
			Owner:		question.Owner,
			Class_id:   oldClass.Class_id,     
			Created_at:   time.Now(),
		}

		result_notification, err := NotificationCollection.InsertOne(ctx, newNotifiaction)
		

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			fmt.Println(classroom_result)
		}
		fmt.Println(result_notification)
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": newQuestion}})


	}
}

func GetAllQuestions() gin.HandlerFunc {
	return func(c *gin.Context) {
		classroom_id := c.Param("classroom_id")
		fmt.Println(classroom_id)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		questions, err := findQuestionsByclassroom_id(ctx, classroom_id)
		
		if err != nil {
			c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Fail to Find Data", Result: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": questions, "count": len(questions)}})

	}
}

func findQuestionsByclassroom_id(ctx context.Context, classroom_id string) ([]interface{}, error) {
	var questions []interface{}
	
	results, err := QeustionCollection.Find(ctx, bson.M{"class_id": classroom_id})
	if err != nil {
		return nil, err
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleQuestion models.Question
		if err = results.Decode(&singleQuestion); err != nil {
			return nil, err
		}
		log.Println(singleQuestion)
		questions = append(questions, singleQuestion)
	}

	return questions, nil
}

func GetQuestionById() gin.HandlerFunc {
	return func(c *gin.Context) {
		classroom_id := c.Param("question_id")
		fmt.Println(classroom_id)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		questions, err := findQuestionsByclassroom_id(ctx, classroom_id)
		if err != nil {
			c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Fail to Find Data", Result: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": questions, "count": len(questions)}})

	}
}

func GetQuestionByuserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		questionId := c.Param("question_id")
		fmt.Println(questionId)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var questions []interface{}
	
		results, err := QeustionCollection.Find(ctx, bson.M{"owner": questionId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
		}
		defer results.Close(ctx)

		for results.Next(ctx) {
			var singleQuestion models.Question
			if err = results.Decode(&singleQuestion); err != nil {
				c.JSON(http.StatusInternalServerError, "error")
			}
			log.Println(singleQuestion)
			questions = append(questions, singleQuestion)
		}
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": questions}})

	}
}

func DeleteQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		questionId := c.Param("questionId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(questionId)

		result, err := QeustionCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, responses.DeleteResponse{Status: http.StatusOK, Message: "Deleted Successfully"})
	}
}

func UpdateQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		questionId := c.Param("question_id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(questionId)

		result, err := QeustionCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, responses.DeleteResponse{Status: http.StatusOK, Message: "Deleted Successfully"})
	}
}
