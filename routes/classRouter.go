package routes

import (
	"exmaple/Backendasktu/controllers"
	middleware "exmaple/Backendasktu/middleware"

	"github.com/gin-gonic/gin"
)

func ClassRoutes(router *gin.Engine) {

	router.Use(middleware.Authentication())
	//routes for version 1 is ready for use
	v1 := router.Group("api/v1")
	{
		v1.GET("/class", controllers.GetAllClassroom())
		v1.POST("/class", controllers.CreateClassroom())
		v1.GET("/class/:classId", controllers.GetClassroom())
		v1.PUT("/class/:classId", controllers.UpdateClassromm())
		v1.DELETE("/class/:classId", controllers.DeleteClassroom())

		//v1.GET("/class/questions", controllers.GetAllQuestions())
		v1.GET("/class/:classId/questions", controllers.GetAllQuestions())
		v1.POST("/class/:classId/question", controllers.CreateQuestion())
		v1.DELETE("/question/:questionId", controllers.DeleteQuestion())

		v1.GET("/class/question/:questionId/answers", controllers.GetAllAnswers())
		//v1.GET("/answer/:answerID", controllers.GetAnswer())
		v1.POST("class/question/:questionId/answer", controllers.CreateAnswer())
	}
	//routes for version 2 is Development in progress
	v2 := router.Group("api/v2")
	{
		v2.GET("/class", controllers.GetAllClassroom())             //get all class
		v2.POST("/class", controllers.CreateClassroom())            // create class with body data
		v2.GET("/class/:classId", controllers.GetClassroom())       //get class by id params
		v2.PUT("/class/:classId", controllers.UpdateClassromm())    //update class by id params and body data
		v2.DELETE("/class/:classId", controllers.DeleteClassroom()) // delete class by id params

		//v2.GET("/class/questions", controllers.GetAllQuestions())
		v2.GET("/class/:classId/questions", controllers.GetAllQuestions()) //get all question by class id
		v2.POST("/class/:classId/question", controllers.CreateQuestion())  //create question by class id and body data
		v2.DELETE("/question/:questionId", controllers.DeleteQuestion())   // delete question by id params

		v2.GET("/class/question/:questionId/answers", controllers.GetAllAnswers()) //get all answers by question id
		//v2.GET("/answer/:answerID", controllers.GetAnswer())
		v2.POST("class/question/:questionId/answer", controllers.CreateAnswer()) //create answer by question id and body data

	}

}
