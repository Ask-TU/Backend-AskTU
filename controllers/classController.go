package controllers

import (
	"exmaple/Backendasktu/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var class = []models.AllClass{
	{Subject_name: "SF230", Class_owner: "Monmanut"},
	{Subject_name: "SF530", Class_owner: "Teetawat"},
}

func GetClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, class)
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
