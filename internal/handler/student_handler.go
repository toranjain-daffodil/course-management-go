package handler

import (
	"assignment/internal/mapper"
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
	students.GET("/unenrolled", h.ListUnenrolledStudents)
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

	response := make([]model.StudentResponse, 0, len(students))
	for _, student := range students {
		response = append(response, mapper.ToStudentResponse(student))
	}
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	student, err := h.studentService.Get(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToStudentResponse(*student))
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var input model.CreateStudentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := h.studentService.Create(c.Request.Context(), input)
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

	student, err := h.studentService.Update(c.Request.Context(), id, input)

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

	if err := h.studentService.Delete(c.Request.Context(), id); err != nil {
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

func (h *StudentHandler) ListUnenrolledStudents(c *gin.Context) {
	students, err := h.studentService.ListUnenrolled(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	response := make([]model.StudentResponse, 0, len(students))
	for _, student := range students {
		response = append(response, mapper.ToStudentResponse(student))
	}
	c.JSON(http.StatusOK, response)
}
