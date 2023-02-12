package routes

import (
	controller "exmaple/Backendasktu/controllers"
	middleware "exmaple/Backendasktu/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {

	router.Use(middleware.Authentication())
	router.GET("/users", controller.GetAllUsers())
	router.GET("/user/:user_id", controller.GetUser())
	router.PUT("/user/:user_id", controller.UpdateUser())
}
