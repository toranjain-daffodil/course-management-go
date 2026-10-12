package router

import (
	"assignment/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	courseHandler *handler.CourseHandler,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	handler.RegisterCourseRoutes(api, courseHandler)

	return r
}