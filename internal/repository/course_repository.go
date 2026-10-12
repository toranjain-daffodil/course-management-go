package repository

import (
	"database/sql"

	"gorm.io/gorm"
)

type CourseRepository struct {
	gormDB *gorm.DB
	sqlDB *sql.DB
}

func NewCourseRepository(
	gormDB *gorm.DB,
	sqlDB *sql.DB,
) *CourseRepository {
	return &CourseRepository{
		gormDB: gormDB,
		sqlDB: sqlDB,
	}
}