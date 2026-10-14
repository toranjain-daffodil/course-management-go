package config

import (
	"assignment/internal/model"
	"database/sql"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() (*model.Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err 
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		return nil, model.ErrMissingConfig
	}

	maxStudents, err := strconv.ParseInt(os.Getenv("MAX_STUDENTS_PER_COURSE"), 10, 64)
	if err != nil {
		return nil, err
	}

	maxCourses, err := strconv.ParseInt(os.Getenv("MAX_COURSES_PER_STUDENT"), 10, 64)
	if err != nil {
		return nil, err
	}
	defaultLimit, err := strconv.ParseInt(os.Getenv("DEFAULT_LIMIT"), 10, 64)
	if err != nil {
		return nil, err
	}
	maxLimit, err := strconv.ParseInt(os.Getenv("MAX_LIMIT"), 10, 64)
	if err != nil {
		return nil, err
	}

	enroll := model.EnrollmentConfig{
		MaxStudentsPerCourse: int(maxStudents),
		MaxCoursesPerStudent: int(maxCourses),
	}
	pagination := model.PaginationConfig{
		DefaultLimit: int(defaultLimit),
		MaxLimit: int(maxLimit),
	}

	return &model.Config{
		PostgresDSN: dsn,
		Enrollment: enroll,
		Pagination: pagination,
	}, nil
}

func ConnectDatabase(dsn string) (*gorm.DB, *sql.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	return gormDB, sqlDB, nil
}