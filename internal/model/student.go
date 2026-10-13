package model

type Student struct {
	Id 					uint 		`json:"id"`
	Name 				string 		`json:"name"`
	Email 				string 		`json:"email" binding:"email"`
	EnrolledCourses 	[]Course 	`gorm:"many2many:course_students" json:"enrolledCourses"`
}

type StudentFilter struct {
	Limit int `form:"limit"`
	Offset int `form:"offset"`
	Search string `form:"search"`
}

type CreateStudentRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateStudentRequest struct {
	Name *string `json:"name"`
	Email *string `json:"email" binding:"email"`
	EnrolledCourses *[]Course `json:"enrolledCourses"`

}