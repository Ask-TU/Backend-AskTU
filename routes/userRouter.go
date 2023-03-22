package routes

import (
	controller "exmaple/Backendasktu/controllers"
	middleware "exmaple/Backendasktu/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {

	router.Use(middleware.Authentication())
	router.GET("/users", controller.GetAllUsers())
	router.GET("/users/:user_id", controller.GetUser())
	router.PUT("/users/:user_id", controller.UpdateUser())
}
