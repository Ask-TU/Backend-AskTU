package routes

import (
	middleware "exmaple/Backendasktu/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ClassRoomRoutes(router *gin.Engine) {

	router.Use(middleware.Authentication())
	//routes for version 1 is ready for use
	v1 := router.Group("api/v1")
	{
		v1.GET("/posts", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Retrieving posts in version 1"})
		})
		v1.GET("/comments", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Retrieving commetns in version 1"})
		})
	}
	//routes for version 2 is Development in progress
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