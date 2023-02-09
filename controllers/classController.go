package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"exmaple/Backendasktu/database"

	"exmaple/Backendasktu/models"

	"go.mongodb.org/mongo-driver/bson"
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
		classId := c.Param("classId")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var class1 models.AllClass
		objId, _ := primitive.ObjectIDFromHex(classId)

		err := classroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&class1)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, class1)
		fmt.Println(class1)
	}
}
func GetAllClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// recordPerPage := 10
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"class_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		result, err := classroomCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allClass []bson.M
		if err = result.All(ctx, &allClass); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allClass[0])

	}
}

func CreateClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var newclass models.AllClass

		if err := c.BindJSON(&newclass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&newclass)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		newUser := models.AllClass{
			ID:           primitive.NewObjectID(),
			Main_id:      newclass.Main_id,
			Subject_name: newclass.Subject_name,
			Class_owner:  newclass.Class_owner,
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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		classId := c.Param("classId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(classId)

		result, err := classroomCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, "deleted success")

	}
}

func UpdateClassromm() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		classId := c.Param("classId")
		var class1 models.AllClass
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(classId)

		if err := c.BindJSON(&class1); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		validationErr := validate.Struct(&class1)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": validationErr.Error()})
			return
		}

		update := bson.M{"main_id": class1.Main_id, "subject_name": class1.Subject_name, "class_owner": class1.Class_owner}
		result, err := classroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error3")
			return
		}
		var updatedClass models.AllClass
		if result.MatchedCount == 1 {
			err := classroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedClass)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "cannot update")
				return
			}
		}

		c.JSON(http.StatusOK, "update success")
	}
}

func GetAllQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var allQuestions []models.Question
		for _, class := range classes {
			allQuestions = append(allQuestions, class.Question...)
		}
		c.JSON(200, gin.H{"questions": allQuestions})
	}
}

func GetQuestionInClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		classID := c.Param("id")
		for _, class := range classes {
			if class.Main_id == classID {
				c.JSON(200, gin.H{"questions": class.Question})
				return
			}
		}
		c.JSON(404, gin.H{"error": fmt.Sprintf("class with ID %s not found", classID)})
	}
}

func GetQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		questionId := c.Param("id")
		for _, class := range classes {
			for _, question := range class.Question {
				if question.Question_id == questionId {
					c.JSON(200, gin.H{"question": question})
					return
				}
			}
		}
		c.JSON(404, gin.H{"error": fmt.Sprintf("question with ID %s not found", questionId)})
	}
}

func CreateQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var question models.Question
		if err := c.ShouldBindJSON(&question); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		for i, class := range classes {
			if class.Main_id == question.Class_id {
				classes[i].Question = append(classes[i].Question, question)
				c.JSON(201, gin.H{"question": question})
				return
			}
		}
		c.JSON(201, gin.H{"question": question})
	}
}
