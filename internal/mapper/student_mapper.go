package mapper

import "assignment/internal/model"

func ToStudentResponse(student model.Student) model.StudentResponse {
	courses := make([]model.StudentCourse, 0, len(student.EnrolledCourses))

	for _, course := range student.EnrolledCourses {
		courses = append(courses, model.StudentCourse{
			Id:              course.Id,
			Title:           course.Title,
			Description:     course.Description,
			DurationInHours: course.DurationInHours,
			StartDate:       course.StartDate,
			Instructor:      course.Instructor,
		})
	}

	return model.StudentResponse{
		Id:              student.Id,
		Name:            student.Name,
		Email:           student.Email,
		EnrolledCourses: courses,
	}
}
