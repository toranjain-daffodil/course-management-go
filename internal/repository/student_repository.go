package repository

import (
	"assignment/internal/model"
	"context"

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
	return []model.Student{}, nil
}

func (repo studentRepositoryGORM) GetByID(c context.Context, id uint) (*model.Student, error) {
	return nil, nil
}

func (repo studentRepositoryGORM) Create(c context.Context, student *model.Student) error {
	return nil
}

func (repo studentRepositoryGORM) Update(c context.Context, student *model.Student) error {
	return nil
}

func (repo studentRepositoryGORM) Delete(c context.Context, id uint) error {
	return nil
}