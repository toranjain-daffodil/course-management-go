package repository

import (
	"assignment/internal/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type courseRepositoryGORM struct {
	db *gorm.DB
}

func NewCourseRepositoryGORM(
	gormDB *gorm.DB,
) *courseRepositoryGORM {
	return &courseRepositoryGORM{
		db: gormDB,
	}
}

func (repo courseRepositoryGORM) List(c context.Context, filter model.CourseFilter) ([]model.Course, error) {
	q := repo.db.WithContext(c).Model(&model.Course{})

	if filter.Search != ""{
		q = q.Where("title ILIKE ?", "%"+filter.Search+"%")
	}

	var courses []model.Course
	err := q.Preload("Students").Order("id ASC").Limit(filter.Limit).Offset(filter.Offset).Find(&courses).Error
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	return courses, nil
}

func (repo courseRepositoryGORM) GetByID(c context.Context, id uint) (*model.Course, error) {
	var course model.Course
	err := repo.db.WithContext(c).Preload("Students").First(&course, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get course %d: %w", id, err)
	}
	return &course, nil
}

func (repo courseRepositoryGORM) Create(c context.Context, course *model.Course) error {
	err := repo.db.WithContext(c).Create(course).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: course title already exists", model.ErrConflict)
	}
	if err != nil {
		return fmt.Errorf("create course: %w", err)
	}
	return nil
}

func (repo courseRepositoryGORM) Update(c context.Context, course *model.Course) error {
	res := repo.
			db.
			WithContext(c).
			Model(course).
			Select(
				"Title", 
				"Description", 
				"DurationInHours", 
				"StartDate", 
				"Instructor",
			).
			Updates(course)
	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: course title already exists", model.ErrConflict)
	}
	if res.Error != nil {
		return fmt.Errorf("update course %d: %w", course.Id, res.Error)
	}
	if res.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (repo courseRepositoryGORM) Delete(c context.Context, id uint) error {
	res := repo.db.WithContext(c).Delete(&model.Course{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete course %d: %w", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (repo courseRepositoryGORM) Enroll(c context.Context, courseID uint, studentID uint) error {
	var course model.Course
	var student model.Student

	if err := repo.db.WithContext(c).First(&course, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound
		}
		return err
	}

	if err := repo.db.WithContext(c).First(&student, studentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrNotFound
		}
		return err
	}

	if err := repo.db.Model(&course).Association("Students").Append(&student); err != nil {
		return err
	}
	return nil
}