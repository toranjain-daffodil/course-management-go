package interfaces

import (
	"assignment/internal/model"
	"context"
)

type CourseRepository interface {
	List(c context.Context, filter model.CourseFilter) ([]model.Course, error)
	GetByID(c context.Context, id uint) (*model.Course, error)
	Create(c context.Context, course *model.Course) error
	Update(c context.Context, course *model.Course) error
	Delete(c context.Context, id uint) error
}