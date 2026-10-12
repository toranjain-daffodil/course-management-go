package router

import (
	"assignment/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	courseHandler *handler.CourseHandler,
	studentHandler *handler.StudentHandler,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	handler.RegisterCourseRoutes(api, courseHandler)
	handler.RegisterStudentRoutes(api, studentHandler)

	return r
}