package service

import (
	"assignment/internal/repository"
)

type CourseService struct {
	courseRepository repository.CourseRepository
}

func NewCourseService(repo *repository.CourseRepository) *CourseService {
	return &CourseService{
		courseRepository: *repo,
	}
}