package model

import "time"

type Course struct {
	Id              uint      `gorm:"primary-key" json:"id"`
	Title           string    `gorm:"size:200;not null;uniqueIndex" json:"title"`
	Description     string    `json:"description"`
	DurationInHours uint      `json:"durationInHours"`
	StartDate       time.Time `json:"startDate"`
	Instructor      string    `json:"instructor"`
	Students        []Student `gorm:"many2many:course_students" json:"students"`
}

type CourseFilter struct {
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Search string `form:"search"`
}

type CreateCourseRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	DurationInHours uint      `json:"durationInHours" binding:"required"`
	StartDate       time.Time `json:"startDate"`
	Instructor      string    `json:"instructor" binding:"required"`
}

type UpdateCourseRequest struct {
	Title           *string    `json:"title"`
	Description     *string    `json:"description"`
	DurationInHours *uint      `json:"durationInHours"`
	StartDate       *time.Time `json:"startDate"`
	Instructor      *string    `json:"instructor"`
}

type EnrollCourseRequest struct {
	StudentID uint `json:"studentID" binding:"required"`
}

type CourseStudent struct {
	Id uint `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CourseResponse struct {
	Id              uint            `json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	DurationInHours uint            `json:"durationInHours"`
	StartDate       time.Time       `json:"startDate"`
	Instructor      string          `json:"instructor"`
	Students        []CourseStudent `json:"students"`
}
