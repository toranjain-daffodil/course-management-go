package service

import (
	"assignment/internal/model"
	"assignment/internal/repository/interfaces"
	"context"
	"strings"
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
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	if filter.Limit > maxLimit {
		filter.Limit = maxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return service.studentRepository.List(c, filter)
}

func (service *studentService) Get(c context.Context, id uint) (*model.Student, error) {
	return service.studentRepository.GetByID(c, id)
}

func (service *studentService) Create(c context.Context, input model.CreateStudentRequest) (*model.Student, error) {
	student := &model.Student{
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(input.Email),
	}

	if err := service.studentRepository.Create(c, student); err != nil {
		return nil, err
	}
	return student, nil
}

func (service *studentService) Update(c context.Context, id uint, input model.UpdateStudentRequest) (*model.Student, error) {
	student, err := service.studentRepository.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		student.Name = name
	}

	if input.Email != nil {
		email := strings.TrimSpace(*input.Email)
		student.Email = email
	}

	if err := service.studentRepository.Update(c, student); err != nil {
		return nil, err
	}
	return student, nil

}

func (service *studentService) Delete(c context.Context, id uint) error {
	return service.studentRepository.Delete(c, id)
}

func (service *studentService) ListUnenrolled(c context.Context) ([]model.Student, error) {
	return service.studentRepository.ListUnenrolled(c)
}
