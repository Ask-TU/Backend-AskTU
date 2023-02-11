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
var answerCollection *mongo.Collection = database.OpenCollection(database.Client, "answers")

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

		update := bson.M{"subject_name": class1.Subject_name, "class_owner": class1.Class_owner}
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

func CreateQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		classId := c.Param("classId")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(classId)

		var class1 models.AllClass

		err := classroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&class1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var newQuestion models.Question

		newQuestion = models.Question{
			ID:          primitive.NewObjectID(),
			Question_id: newQuestion.Question_id,
			Content:     newQuestion.Content,
			Owner:       newQuestion.Owner,
		}

		if err := c.BindJSON(&newQuestion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		class1.Question = append(class1.Question, newQuestion)

		update := bson.M{"$set": bson.M{"question": class1.Question}}
		_, err = classroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, class1)
	}
}

func DeleteQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		questionId := c.Param("questionId")

		objId, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := classroomCollection.UpdateOne(ctx, bson.M{"question._id": objId}, bson.M{"$pull": bson.M{"question": bson.M{"_id": objId}}})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
	}
}

func GetQuestionsByClassID() gin.HandlerFunc {
	return func(c *gin.Context) {
		classID := c.Param("classId")

		objID, err := primitive.ObjectIDFromHex(classID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
			return
		}

		var class models.AllClass
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = classroomCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&class)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, class.Question)
	}
}

func GetAllQuestions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var classrooms []models.AllClass

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := classroomCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := cursor.All(ctx, &classrooms); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var questions []models.Question
		for _, class := range classrooms {
			questions = append(questions, class.Question...)
		}

		c.JSON(http.StatusOK, questions)
	}
}

func CreateAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var newAnswer models.Answer

		if err := c.BindJSON(&newAnswer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newRespon := models.Answer{
			ID:          primitive.NewObjectID(),
			Content:     newAnswer.Content,
			Owner:       newAnswer.Owner,
			Question_id: newAnswer.Question_id,
			Answer_id:   newAnswer.Answer_id,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}

		result, err := answerCollection.InsertOne(ctx, newRespon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result)
		c.JSON(http.StatusCreated, "success")

	}
}

func GetAllAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		allans, err := strconv.Atoi(c.Query("allans"))
		if err != nil || allans < 1 {
			allans = 10
		}

		face, err1 := strconv.Atoi(c.Query("face"))
		if err1 != nil || face < 1 {
			face = 1
		}

		startIndex := (face - 1) * allans
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"class_items", bson.D{{"$slice", []interface{}{"$data", startIndex, allans}}}},
			}}}

		result, err := answerCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allanswers []bson.M
		if err = result.All(ctx, &allanswers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allanswers[0])
	}
}

func GetAnswerInQuestion() gin.HandlerFunc {
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
