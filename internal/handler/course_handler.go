package handler

import (
	"assignment/internal/mapper"
	"assignment/internal/model"
	"assignment/internal/service/interfaces"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService interfaces.CourseService
}

func NewCourseHandler(service interfaces.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: service,
	}
}

func RegisterCourseRoutes(
	api *gin.RouterGroup,
	h *CourseHandler,
) {
	courses := api.Group("/courses")
	courses.GET("", h.ListCourses)
	courses.GET("/:id", h.GetCourse)
	courses.POST("", h.CreateCourse)
	courses.PUT("/:id", h.UpdateCourse)
	courses.DELETE("/:id", h.DeleteCourse)
	courses.POST("/:id/enroll", h.EnrollCourse)
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	var filter model.CourseFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courses, err := h.courseService.List(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response := make([]model.CourseResponse, 0, len(courses))

	for _, course := range courses {
		response = append(response, mapper.ToCourseResponse(course))
	}
	c.JSON(http.StatusOK, response)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	course, err := h.courseService.Get(c, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := mapper.ToCourseResponse(*course)
	c.JSON(http.StatusOK, response)

}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var input model.CreateCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseService.Create(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var input model.UpdateCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.courseService.Update(c, id, input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.courseService.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CourseHandler) EnrollCourse(c *gin.Context) {
	courseID, ok := parseID(c)
	if !ok {
		return
	}
	
	var input model.EnrollCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	if err := h.courseService.Enroll(c, courseID ,input.StudentID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id), true
}

func (h *CourseHandler) handleError(c *gin.Context, err error) {
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