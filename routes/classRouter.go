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
		v1.GET("/class", controllers.GetAllClassrooms())
		v1.POST("/class", controllers.CreateClassroom())
		v1.GET("/class/:classId", controllers.GetClassroom())
		v1.PUT("/class/:classId", controllers.UpdateClassroom())
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
		//router.GET("/notifications", controllers.GetAllNotifications)
		//router.GET("/notifications/:notification_id", controllers.GetNotification)
		//router.POST("/notifications", controllers.CreateNotification)
		//router.DELETE("/notifications/:notification_id", controllers.DeleteNotification)

		// Question endpoints
		v2.GET("/classrooms/:classroom_id/questions", controllers.GetAllQuestions())
		//v2.GET("/classrooms/:classroom_id/questions/:question_id", controllers.GetQuestionById())
		v2.POST("/classrooms/:classroom_id/questions", controllers.CreateQuestion())
		//v2.PUT("/classrooms/:classroom_id/questions/:question_id", controllers.UpdateQuestion())
		///router.DELETE("/questions/:question_id", controllers.DeleteQuestion)

		// answer endpoints
		//router.GET("/questions/:question_id/answers", controllers.GetAllanswers)
		v2.POST("/classrooms/questions/:question_id/answers", controllers.CreateAnswer())
		//router.PUT("/answers/:answer_id", controllers.Updateanswer)
		//router.DELETE("/answers/:answer_id", controllers.DeleteComment)

		// Classroom endpoints
		v2.GET("/classrooms", controllers.GetAllClassrooms())
		v2.GET("/classrooms/:classroom_id", controllers.GetClassroom())
		v2.POST("/classrooms", controllers.CreateClassroom())
		v2.PUT("/classrooms/:classroom_id", controllers.UpdateClassroom())
		v2.DELETE("/classrooms/:classroom_id", controllers.DeleteClassroom())
		//
		// Classroom Members endpoints
		//router.GET("/classrooms/:classroom_id/members", controllers.GetClassroomMembers)
		v2.POST("/classrooms/join/:classroom_id/:member_id", controllers.JoinClasrooms())
		v2.DELETE("/classrooms/delete/:classroom_id/:member_id", controllers.RemoveClassroomMember())

		// Search endpoints
		//router.GET("/questions/search", controllers.SearchQuestions)
		////router.GET("/comments/search", controllers.SearchComments)
		////router.GET("/classrooms/search", controllers.SearchClassrooms)
		//router.GET("/search", controllers.GlobalSearch)

	}

}
