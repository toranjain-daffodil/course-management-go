package handler

import (
	"assignment/internal/model"
	"assignment/internal/service/interfaces"
	"errors"
	"net/http"

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
	students := api.Group("/students")
	students.GET("", h.ListStudent)
	students.GET("/:id", h.GetStudent)
	students.POST("", h.CreateStudent)
	students.PUT("/:id", h.UpdateStudent)
	students.DELETE("/:id", h.DeleteStudent)
}

func (h *StudentHandler) ListStudent(c *gin.Context) {
	var filter model.StudentFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	students, err := h.studentService.List(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	course, err := h.studentService.Get(c, id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var input model.CreateStudentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := h.studentService.Create(c, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var input model.UpdateStudentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := h.studentService.Update(c, id, input)
	
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.studentService.Delete(c, id); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *StudentHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
	case errors.Is(err, model.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}