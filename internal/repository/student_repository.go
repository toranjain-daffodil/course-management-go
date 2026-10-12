package repository

import (
	"assignment/internal/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type studentRepositoryGORM struct {
	db *gorm.DB
}

func NewStudentRepositoryGORM(db *gorm.DB) *studentRepositoryGORM {
	return &studentRepositoryGORM{
		db: db,
	}
}

func (repo studentRepositoryGORM) List(c context.Context, filter model.StudentFilter) ([]model.Student, error) {
	q := repo.db.WithContext(c).Model(&model.Student{})

	if filter.Search != "" {
		q = q.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	var students []model.Student
	res := q.Order("id ASC").Limit(filter.Limit).Offset(filter.Offset).Find(&students)
	if res.Error != nil {
		return nil, fmt.Errorf("list students: %w", res.Error)
	}
	return students, nil
}

func (repo studentRepositoryGORM) GetByID(c context.Context, id uint) (*model.Student, error) {
	var student model.Student
	res := repo.db.WithContext(c).First(&student, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	}
	if res.Error != nil {
		return nil, fmt.Errorf("get student %d: %w", id, res.Error)
	}
	return &student, nil
}

func (repo studentRepositoryGORM) Create(c context.Context, student *model.Student) error {
	res := repo.db.WithContext(c).Create(student)
	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: student already exists", model.ErrConflict)
	}
	if res.Error != nil {
		return fmt.Errorf("create student: %w", res.Error)
	}
	return nil
}

func (repo studentRepositoryGORM) Update(c context.Context, student *model.Student) error {
	res := repo.db.WithContext(c).
			Model(student).
			Select(
				"Name",
				"Email",
				"EnrolledCourses",
			).
			Updates(student)
	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: student already exists", model.ErrConflict)
	}
	if res.Error != nil {
		return fmt.Errorf("update student %d: %w", student.Id, res.Error)
	}
	if res.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (repo studentRepositoryGORM) Delete(c context.Context, id uint) error {
	res := repo.db.WithContext(c).Delete(&model.Student{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete student %d: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}