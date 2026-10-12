package service

import (
	"assignment/internal/model"
	"assignment/internal/repository/interfaces"
	"context"
)

type studentService struct {
	studentRepository interfaces.StudentRepository
}

func NewStudentService(repo interfaces.StudentRepository) *studentService {
	return &studentService{
		studentRepository: repo,
	}
}

func (service *studentService) List(c context.Context, filter model.StudentFilter) ([]model.Student, error) {
	return nil, nil
}

func (service *studentService) Get(c context.Context, id uint) (*model.Student, error) {
	return nil, nil
}

func (service *studentService) Create(c context.Context, input model.CreateStudentRequest) (*model.Student, error) {
	return nil, nil
}

func (service *studentService) Update(c context.Context, id uint, input model.UpdateStudentRequest) (*model.Student, error) {
	return nil, nil
}

func (service *studentService) Delete(c context.Context, id uint) error {
	return nil
}