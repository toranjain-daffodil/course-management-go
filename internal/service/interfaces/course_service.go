package interfaces

import (
	"assignment/internal/model"
	"context"
)

type CourseService interface {
	List(c context.Context, filter model.CourseFilter) ([]model.Course, error)
	Get(c context.Context, id uint) (*model.Course, error)
	Create(c context.Context, input model.CreateCourseRequest) (*model.Course, error)
	Update(c context.Context, id uint, input model.UpdateCourseRequest) (*model.Course, error)
	Delete(c context.Context, id uint) error
}