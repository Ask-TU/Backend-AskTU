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

var ClassroomCollection *mongo.Collection = database.OpenCollection(database.Client, "classrooms")
var NotificationCollection *mongo.Collection = database.OpenCollection(database.Client, "notifications")

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

		count, err := ClassroomCollection.CountDocuments(ctx, bson.M{"subject_name": class.Subject_name})
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
			Section: 	class.Section,
		}
		result, err := ClassroomCollection.InsertOne(ctx, newClass)
				if err != nil {
					c.JSON(http.StatusInternalServerError, "error")
					return
		}
		
		newNotifiaction := models.Notification{
			ID:           primitive.NewObjectID(),
			Content:	  "You have created a new class",
			Owner:		newClass.Class_owner,
			Class_id:   newClass.Class_id,     
			Created_at:   time.Now(),
		}

		result_notification, err := NotificationCollection.InsertOne(ctx, newNotifiaction)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}
		fmt.Print(result, result_notification)
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": newClass}})

	}
}

func GetClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		classroom_id := c.Param("classroom_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var oldClass models.Classrooms
		objId, _ := primitive.ObjectIDFromHex(classroom_id)
		fmt.Print(objId)
		err := ClassroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&oldClass)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, oldClass)
		fmt.Println(oldClass)
	}
}

func GetAllClassrooms() gin.HandlerFunc {
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

		result, err := ClassroomCollection.Aggregate(ctx, mongo.Pipeline{
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
		classroom_id := c.Param("classroom_id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(classroom_id)

		result, err := ClassroomCollection.DeleteOne(ctx, bson.M{"_id": objId})

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

func UpdateClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		classroom_id := c.Param("classroom_id")

		var oldClass models.Classrooms
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(classroom_id)

		if err := c.BindJSON(&oldClass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		validationErr := validate.Struct(&oldClass)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": validationErr.Error()})
			return
		}

		update := bson.M{
			"class_id":     oldClass.Class_id,
			"subject_name": oldClass.Subject_name,
			"class_owner":  oldClass.Class_owner,
			"section":      oldClass.Section,
			"members":      oldClass.Members,
			"questions":    oldClass.Questions,
			
		}
		result, err := ClassroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, "error3")
			return
		}
		var updatedClass models.Classrooms

		if result.MatchedCount == 1 {
			err := ClassroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedClass)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "cannot update")
				return
			}
		}

		c.JSON(http.StatusOK, responses.UpdateResponse{Status: http.StatusOK, Message: "Update Successfully"})
	}
}

func JoinClasrooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroom_id := c.Param("classroom_id")
		member_id := c.Param("member_id")

		objId, err := primitive.ObjectIDFromHex(classroom_id)

		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		}
		//2
		var oldClass models.Classrooms
		err = ClassroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&oldClass)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newMember := []string{member_id}
		mergeMember := append(oldClass.Members, newMember...)
		update := bson.M{
			"subject_name": oldClass.Subject_name,
			"class_owner":  oldClass.Class_owner,
			"created_at":   oldClass.Created_at,
			"updated_at":   time.Now(),
			"questions":    oldClass.Questions,
			"members":      mergeMember,
		}

		result, err := ClassroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			fmt.Println(result)
		}
		//3
		objUserId, err := primitive.ObjectIDFromHex(member_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		}

		var oldUser models.User
		err = Usercollection.FindOne(ctx, bson.M{"_id": objUserId}).Decode(&oldUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newUserInfomation := []string{classroom_id}
		mergeUserInfomation := append(oldUser.Classrooms, newUserInfomation...)
		updateInformation := bson.M{
			"first_name": oldUser.First_name,
			"last_name":  oldUser.Last_name,
			"nick_name":  oldUser.Nick_name,
			"email":      oldUser.Email,
			"phone":      oldUser.Phone,
			"password":   oldUser.Password,
			"created_at": oldUser.Created_at,
			"updated_at": time.Now(),
			"user_id":    oldUser.User_id,
			"student_id": oldUser.Student_id,
			"classrooms": mergeUserInfomation,
		}

		result, err = Usercollection.UpdateOne(ctx, bson.M{"_id": objUserId}, bson.M{"$set": updateInformation})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			fmt.Println(result)
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": "Join Classrooms Successfully"}})
	}
}

func RemoveClassroomMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroom_id := c.Param("classroom_id")
		member_id := c.Param("member_id")

		objId, err := primitive.ObjectIDFromHex(classroom_id)

		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		}
		//2 classrooms
		var oldClass models.Classrooms
		err = ClassroomCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&oldClass)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		valueToRemove := member_id
		found := false
		index := -1
		for i, v := range oldClass.Members {
			if v == valueToRemove {
				found = true
				index = i
				break
			}
		}
		if found {
			oldClass.Members = append(oldClass.Members[:index], oldClass.Members[index+1:]...)
		}

		update := bson.M{
			"subject_name": oldClass.Subject_name,
			"class_owner":  oldClass.Class_owner,
			"created_at":   oldClass.Created_at,
			"updated_at":   time.Now(),
			"questions":    oldClass.Questions,
			"members":      oldClass.Members,
		}

		result, err := ClassroomCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			fmt.Println(result)
		}
		//3
		objUserId, err := primitive.ObjectIDFromHex(member_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		}

		var oldUser models.User
		err = Usercollection.FindOne(ctx, bson.M{"_id": objUserId}).Decode(&oldUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		valueToRemove = classroom_id
		found = false
		index = -1
		for i, v := range oldUser.Classrooms {
			if v == valueToRemove {
				found = true
				index = i
				break
			}
		}
		if found {
			oldUser.Classrooms = append(oldUser.Classrooms[:index], oldUser.Classrooms[index+1:]...)
		}
		updateInformation := bson.M{
			"first_name": oldUser.First_name,
			"last_name":  oldUser.Last_name,
			"nick_name":  oldUser.Nick_name,
			"email":      oldUser.Email,
			"phone":      oldUser.Phone,
			"password":   oldUser.Password,
			"created_at": oldUser.Created_at,
			"updated_at": time.Now(),
			"user_id":    oldUser.User_id,
			"student_id": oldUser.Student_id,
			"classrooms": oldUser.Classrooms,
		}

		result, err = Usercollection.UpdateOne(ctx, bson.M{"_id": objUserId}, bson.M{"$set": updateInformation})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Result: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			fmt.Println(result)
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "Successfully", Result: map[string]interface{}{"data": "Delete Member from Classrooms Successfully"}})
	}
}

