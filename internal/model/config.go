package model

type Config struct {
	PostgresDSN string
	Enrollment EnrollmentConfig
	Pagination PaginationConfig
}

type EnrollmentConfig struct {
	MaxStudentsPerCourse int
	MaxCoursesPerStudent int
}

type PaginationConfig struct {
	DefaultLimit int
	MaxLimit int
}