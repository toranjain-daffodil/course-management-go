package repository

import (
	"assignment/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type courseRepositorySQL struct {
	db     *sql.DB
	enroll model.EnrollmentConfig
}

func NewCourseRepositorySQL(
	db *sql.DB,
	enroll model.EnrollmentConfig,
) *courseRepositorySQL {
	return &courseRepositorySQL{
		db:     db,
		enroll: enroll,
	}
}

func (r courseRepositorySQL) List(c context.Context, filter model.CourseFilter) ([]model.Course, error) {
	rows, err := r.db.QueryContext(c, `
		SELECT id, title, description, duration_in_hours, start_date, instructor
		FROM courses
		WHERE ($1 = '' OR title ILIKE '%' || $1 || '%')
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`, filter.Search, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	defer rows.Close()

	courses, err := scanCourses(rows)
	if err != nil {
		return nil, err
	}

	if err := r.loadStudents(c, courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r courseRepositorySQL) GetByID(c context.Context, id uint) (*model.Course, error) {
	var course model.Course

	err := r.db.QueryRowContext(c, `
		SELECT id, title, description, duration_in_hours, start_date, instructor
		FROM courses
		WHERE id = $1
	`, id).Scan(
		&course.Id,
		&course.Title,
		&course.Description,
		&course.DurationInHours,
		&course.StartDate,
		&course.Instructor,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get course %d: %w", id, err)
	}

	courses := []model.Course{course}
	if err := r.loadStudents(c, courses); err != nil {
		return nil, err
	}
	return &courses[0], nil
}

func (r courseRepositorySQL) Create(c context.Context, course *model.Course) error {
	err := r.db.QueryRowContext(c, `
		INSERT INTO courses (
			title, description, duration_in_hours, start_date, instructor
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		course.Title,
		course.Description,
		course.DurationInHours,
		course.StartDate,
		course.Instructor,
	).Scan(&course.Id)

	if isUniqueViolation(err) {
		return fmt.Errorf("%w: course title already exists", model.ErrConflict)
	}
	if err != nil {
		return fmt.Errorf("create course: %w", err)
	}
	return nil
}

func (r courseRepositorySQL) Update(c context.Context, course *model.Course) error {
	result, err := r.db.ExecContext(c, `
		UPDATE courses
		SET
			title = $1,
			description = $2,
			duration_in_hours = $3,
			start_date = $4,
			instructor = $5
		WHERE id = $6
	`,
		course.Title,
		course.Description,
		course.DurationInHours,
		course.StartDate,
		course.Instructor,
		course.Id,
	)

	if isUniqueViolation(err) {
		return fmt.Errorf("%w: course title already exists", model.ErrConflict)
	}
	if err != nil {
		return fmt.Errorf("update course %d: %w", course.Id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update course %d: %w", course.Id, err)
	}
	if affected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r courseRepositorySQL) Delete(c context.Context, id uint) error {
	result, err := r.db.ExecContext(c, `
		DELETE FROM courses WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete course %d: %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete course %d: %w", id, err)
	}
	if affected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r courseRepositorySQL) Enroll(
	c context.Context,
	courseID uint,
	studentID uint,
) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("begin enrollment transaction: %w", err)
	}
	defer tx.Rollback()

	var id uint
	err = tx.QueryRowContext(c,
		`SELECT id FROM courses WHERE id = $1 FOR UPDATE`,
		courseID,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("get course for enrollment: %w", err)
	}

	err = tx.QueryRowContext(c,
		`SELECT id FROM students WHERE id = $1 FOR UPDATE`,
		studentID,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("get student for enrollment: %w", err)
	}

	var existing int
	if err := tx.QueryRowContext(c, `
		SELECT COUNT(*)
		FROM course_students
		WHERE course_id = $1 AND student_id = $2
	`, courseID, studentID).Scan(&existing); err != nil {
		return fmt.Errorf("check enrollment: %w", err)
	}
	if existing > 0 {
		return fmt.Errorf("%w: student is already enrolled in this course", model.ErrConflict)
	}

	var courseCount int
	if err := tx.QueryRowContext(c, `
		SELECT COUNT(*)
		FROM course_students
		WHERE course_id = $1
	`, courseID).Scan(&courseCount); err != nil {
		return fmt.Errorf("count course enrollments: %w", err)
	}
	if courseCount >= r.enroll.MaxStudentsPerCourse {
		return fmt.Errorf(
			"%w: a course cannot have more than %d students",
			model.ErrValidation,
			r.enroll.MaxStudentsPerCourse,
		)
	}

	var studentCount int
	if err := tx.QueryRowContext(c, `
		SELECT COUNT(*)
		FROM course_students
		WHERE student_id = $1
	`, studentID).Scan(&studentCount); err != nil {
		return fmt.Errorf("count student enrollments: %w", err)
	}
	if studentCount >= r.enroll.MaxCoursesPerStudent {
		return fmt.Errorf(
			"%w: a student cannot enroll in more than %d courses",
			model.ErrValidation,
			r.enroll.MaxCoursesPerStudent,
		)
	}

	if _, err := tx.ExecContext(c, `
		INSERT INTO course_students (course_id, student_id)
		VALUES ($1, $2)
	`, courseID, studentID); err != nil {
		return fmt.Errorf("enroll student: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit enrollment: %w", err)
	}
	return nil
}

func (r courseRepositorySQL) TopEnrolled(
	c context.Context,
	limit int,
) ([]model.Course, error) {
	rows, err := r.db.QueryContext(c, `
		SELECT
			c.id, c.title, c.description, c.duration_in_hours,
			c.start_date, c.instructor
		FROM courses c
		LEFT JOIN course_students cs ON cs.course_id = c.id
		GROUP BY c.id
		ORDER BY COUNT(cs.student_id) DESC, c.id ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list top enrolled courses: %w", err)
	}
	defer rows.Close()

	courses, err := scanCourses(rows)
	if err != nil {
		return nil, err
	}

	if err := r.loadStudents(c, courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r courseRepositorySQL) Summary(c context.Context) (model.CourseSummary, error) {
	var summary model.CourseSummary

	err := r.db.QueryRowContext(c, `
		SELECT
			(SELECT COUNT(*) FROM courses),
			(SELECT COUNT(*) FROM students),
			COALESCE((SELECT AVG(duration_in_hours)::float8 FROM courses), 0)
	`).Scan(
		&summary.TotalCourses,
		&summary.TotalStudents,
		&summary.AverageDurationHours,
	)
	if err != nil {
		return summary, fmt.Errorf("get course summary: %w", err)
	}
	return summary, nil
}

func (r courseRepositorySQL) loadStudents(
	c context.Context,
	courses []model.Course,
) error {
	if len(courses) == 0 {
		return nil
	}

	args := make([]any, len(courses))
	placeholders := make([]string, len(courses))
	studentsByCourse := make(map[uint][]model.Student, len(courses))

	for i, course := range courses {
		args[i] = course.Id
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		studentsByCourse[course.Id] = []model.Student{}
	}

	rows, err := r.db.QueryContext(c, fmt.Sprintf(`
		SELECT
			cs.course_id,
			s.id, s.name, s.email
		FROM course_students cs
		JOIN students s ON s.id = cs.student_id
		WHERE cs.course_id IN (%s)
		ORDER BY cs.course_id ASC, s.id ASC
	`, strings.Join(placeholders, ",")), args...)
	if err != nil {
		return fmt.Errorf("load course students: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var courseID uint
		var student model.Student

		if err := rows.Scan(
			&courseID,
			&student.Id,
			&student.Name,
			&student.Email,
		); err != nil {
			return fmt.Errorf("scan course student: %w", err)
		}

		studentsByCourse[courseID] = append(studentsByCourse[courseID], student)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate course students: %w", err)
	}

	for i := range courses {
		courses[i].Students = studentsByCourse[courses[i].Id]
	}
	return nil
}

func scanCourses(rows *sql.Rows) ([]model.Course, error) {
	courses := make([]model.Course, 0)

	for rows.Next() {
		var course model.Course

		if err := rows.Scan(
			&course.Id,
			&course.Title,
			&course.Description,
			&course.DurationInHours,
			&course.StartDate,
			&course.Instructor,
		); err != nil {
			return nil, fmt.Errorf("scan course: %w", err)
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate courses: %w", err)
	}
	return courses, nil
}