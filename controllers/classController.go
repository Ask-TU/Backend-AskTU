package controllers

import (
	"exmaple/Backendasktu/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

		var new_class models.AllClass

		if err := c.ShouldBindJSON(&new_class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		class = append(class, new_class)
		c.JSON(http.StatusCreated, class)
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
