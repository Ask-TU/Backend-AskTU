package main

import (
	routes "exmaple/Backendasktu/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.Use(cors.Default())

	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.ClassRoomRoutes(router)

	router.Run(":" + port)
}
