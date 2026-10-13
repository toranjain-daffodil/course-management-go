package service

import (
	"assignment/internal/model"
	"assignment/internal/repository/interfaces"
	"context"
	"fmt"
	"strings"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type courseService struct {
	courseRepository interfaces.CourseRepository
}

func NewCourseService(repo interfaces.CourseRepository) *courseService {
	return &courseService{
		courseRepository: repo,
	}
}

func (service courseService) List(c context.Context, filter model.CourseFilter) ([]model.Course, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	}
	if filter.Limit > maxLimit {
		filter.Limit = maxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return service.courseRepository.List(c, filter)
}

func (service courseService) Get(c context.Context, id uint) (*model.Course, error) {
	return service.courseRepository.GetByID(c, id)
}

func (service courseService) Create(c context.Context, input model.CreateCourseRequest) (*model.Course, error) {
	title := strings.TrimSpace(input.Title)
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	course := &model.Course{
		Title:           title,
		Description:     strings.TrimSpace(input.Description),
		Instructor:      strings.TrimSpace(input.Instructor),
		DurationInHours: input.DurationInHours,
		StartDate:       input.StartDate,
	}
	if err := service.courseRepository.Create(c, course); err != nil {
		return nil, err
	}
	return course, nil
}

func (service courseService) Update(c context.Context, id uint, input model.UpdateCourseRequest) (*model.Course, error) {
	course, err := service.courseRepository.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if err := validateTitle(title); err != nil {
			return nil, err
		}
		course.Title = title
	}

	if input.Description != nil {
		course.Description = strings.TrimSpace(*input.Description)
	}

	if input.Instructor != nil {
		course.Instructor = strings.TrimSpace(*input.Instructor)
	}

	if input.DurationInHours != nil {
		if err := validateDuration(*input.DurationInHours); err != nil {
			return nil, err
		}
		course.DurationInHours = *input.DurationInHours
	}

	if input.StartDate != nil {
		course.StartDate = *input.StartDate
	}

	if err := service.courseRepository.Update(c, course); err != nil {
		return nil, err
	}
	return course, nil
}

func (service courseService) Delete(c context.Context, id uint) error {
	return service.courseRepository.Delete(c, id)
}

func (service courseService) Enroll(c context.Context, courseID uint, studentID uint) error {
	return service.courseRepository.Enroll(c, courseID, studentID)
}

func (service courseService) TopEnrolled(c context.Context) ([]model.Course, error) {
	return service.courseRepository.TopEnrolled(c, 5)
}

func (service courseService) Summary(c context.Context) (model.CourseSummary, error) {
	return service.courseRepository.Summary(c)
}

func validateTitle(title string) error {
	if title == "" {
		return fmt.Errorf("%w: title is required", model.ErrValidation)
	}

	if len(title) > 200 {
		return fmt.Errorf("%w: title must be at most 200 characters", model.ErrValidation)
	}

	return nil
}

func validateDuration(duration uint) error {
	if duration <= 0 {
		return fmt.Errorf("%w: duration must be greater than 0", model.ErrValidation)
	}
	return nil
}
