package routes

import (
	middleware "exmaple/Backendasktu/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ClassroomRoutes(router *gin.Engine) {

	router.Use(middleware.Authentication())

	v1 := router.Group("api/v1")
	{
		v1.GET("/posts", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Retrieving posts in version 1"})
		})
		v1.GET("/comments", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Retrieving commetns in version 1"})
		})
	}

	v2 := router.Group("api/v2")
	{
		v2.GET("/posts", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Retrieving posts in version 2"})
		})
		v2.GET("/comments", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Retrieving commetns in version 2"})
		})
	}

}
