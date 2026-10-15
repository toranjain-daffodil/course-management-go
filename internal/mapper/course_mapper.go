package mapper

import "assignment/internal/model"

func ToCourseResponse(course model.Course) model.CourseResponse {
	students := make([]model.CourseStudent, 0, len(course.Students))

	for _, student := range course.Students {
		students = append(students, model.CourseStudent{
			Id:    student.Id,
			Name:  student.Name,
			Email: student.Email,
		})
	}

	return model.CourseResponse{
		Id:              course.Id,
		Title:           course.Title,
		Description:     course.Description,
		DurationInHours: course.DurationInHours,
		StartDate:       course.StartDate,
		Instructor:      course.Instructor,
		Students:        students,
	}
}
