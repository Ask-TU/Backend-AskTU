package routes

import (
	controller "exmaple/Backendasktu/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/auth/signup", controller.SignUp())
	incomingRoutes.POST("/auth/login", controller.Login())
}
