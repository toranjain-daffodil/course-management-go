package interfaces

import (
	"assignment/internal/model"
	"context"
)

type StudentService interface {
	List(c context.Context, filter model.StudentFilter) ([]model.Student, error)
	Get(c context.Context, id uint) (*model.Student, error)
	Create(c context.Context, input model.CreateStudentRequest) (*model.Student, error)
	Update(c context.Context, id uint, input model.UpdateStudentRequest) (*model.Student, error)
	Delete(c context.Context, id uint) error
	ListUnenrolled(c context.Context) ([]model.Student, error)
}
