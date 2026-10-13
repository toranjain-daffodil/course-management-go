package interfaces

import (
	"assignment/internal/model"
	"context"
)

type StudentRepository interface {
	List(c context.Context, filter model.StudentFilter) ([]model.Student, error)
	GetByID(c context.Context, id uint) (*model.Student, error)
	Create(c context.Context, student *model.Student) error
	Update(c context.Context, student *model.Student) error
	Delete(c context.Context, id uint) error
	ListUnenrolled(c context.Context) ([]model.Student, error)
}
