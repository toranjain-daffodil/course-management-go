package handler

import (
	"assignment/internal/service/interfaces"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService interfaces.StudentService
}

func NewStudentHandler(service interfaces.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: service,
	}
}

func RegisterStudentRoutes(
	api *gin.RouterGroup,
	h *StudentHandler,
) {
	students := api.Group("/student")
	students.GET("", func(ctx *gin.Context) {})
}