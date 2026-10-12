package model

type student struct {
	Id 				uint 		`json:"id" binding:"required"`
	Name 			string 		`json:"name" binding:"required"`
	Email 			string 		`json:"email" binding:"required, email"`
	EnrolledCourse 	[]Course 	`json:"enrolledCourse" binding:"required"`
}