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
	"exmaple/Backendasktu/helpers"

	"exmaple/Backendasktu/models"

	"exmaple/Backendasktu/responses"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var classroomCollection *mongo.Collection = database.OpenCollection(database.Client, "classrooms")

// Add class id to user profile  when user create class and add user id to class member
func CreateClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var class models.Classrooms

		if err := c.BindJSON(&class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&class)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := classroomCollection.CountDocuments(ctx, bson.M{"subject_name": class.Subject_name})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the Class room"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The Class has already exists"})
			return
		}

		newClass := models.Classrooms{
			ID:           primitive.NewObjectID(),
			Class_id:     helpers.GenerateID(),
			Subject_name: class.Subject_name,
			Class_owner:  class.Class_owner,
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
			Questions:    class.Questions,
			Members:      class.Members,
		}

		result, err := classroomCollection.InsertOne(ctx, newClass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result)
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": newClass}})

	}
}

func GetClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		classId := c.Param("classId")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var oldClass models.Classrooms
		objId, _ := primitive.ObjectIDFromHex(classId)

		err := classroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&oldClass)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, oldClass)
		fmt.Println(oldClass)
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

		c.JSON(http.StatusOK, responses.DeleteResponse{Status: http.StatusOK, Message: "Deleted Successfully"})
	}
}

func UpdateClassromm() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		classId := c.Param("classId")
		var oldClass models.Classrooms
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(classId)

		if err := c.BindJSON(&oldClass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		validationErr := validate.Struct(&oldClass)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": validationErr.Error()})
			return
		}

		update := bson.M{"class_id": oldClass.Class_id, "subject_name": oldClass.Subject_name, "class_owner": oldClass.Class_owner, "members": oldClass.Members, "questions": oldClass.Questions}
		result, err := classroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error3")
			return
		}
		var updatedClass models.Classrooms
		if result.MatchedCount == 1 {
			err := classroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedClass)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "cannot update")
				return
			}
		}

		c.JSON(http.StatusOK, responses.UpdateResponse{Status: http.StatusOK, Message: "Update Successfully"})
	}
}

// Add class id to user profile  when user create class and add user id to class member
func CreateFirstClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var class models.Classrooms

		if err := c.BindJSON(&class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&class)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := classroomCollection.CountDocuments(ctx, bson.M{"subject_name": class.Subject_name})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the Class room"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The Class has already exists"})
			return
		}

		newClass := models.Classrooms{
			ID:           primitive.NewObjectID(),
			Class_id:     helpers.GenerateID(),
			Subject_name: class.Subject_name,
			Class_owner:  class.Class_owner,
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
			Questions:    class.Questions,
			Members:      class.Members,
		}

		result, err := classroomCollection.InsertOne(ctx, newClass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result)
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": newClass}})

	}
}

// Add class id to user profile and add user id to class profile
func JoinClassromm() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		classId := c.Param("classId")
		var oldClass models.Classrooms
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(classId)

		if err := c.BindJSON(&oldClass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		validationErr := validate.Struct(&oldClass)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": validationErr.Error()})
			return
		}

		update := bson.M{"class_id": oldClass.Class_id, "subject_name": oldClass.Subject_name, "class_owner": oldClass.Class_owner, "members": oldClass.Members, "questions": oldClass.Questions}
		result, err := classroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error3")
			return
		}
		var updatedClass models.Classrooms
		if result.MatchedCount == 1 {
			err := classroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedClass)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "cannot update")
				return
			}
		}

		c.JSON(http.StatusOK, responses.UpdateResponse{Status: http.StatusOK, Message: "Update Successfully"})
	}
}
