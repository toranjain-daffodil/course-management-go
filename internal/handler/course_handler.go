package handler

import (
	"assignment/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(service *service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: *service,
	}
}

func RegisterCourseRoutes(
	api *gin.RouterGroup,
	h *CourseHandler,
) {
	courses := api.Group("/course")
	courses.GET("/health", h.Health)
}

func (handler *CourseHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}