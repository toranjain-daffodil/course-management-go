package repository

import (
	"assignment/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type studentRepositorySQL struct {
	db *sql.DB
}

func NewStudentRepositorySQL(db *sql.DB) *studentRepositorySQL {
	return &studentRepositorySQL{
		db: db,
	}
}

func (r studentRepositorySQL) List(c context.Context, filter model.StudentFilter) ([]model.Student, error) {
	rows, err := r.db.QueryContext(c, `
		SELECT id, name, email
		FROM students
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`, filter.Search, filter.Limit, filter.Offset)

	if err != nil {
		return nil, fmt.Errorf("list students: %w", err)
	}
	defer rows.Close()
	
	students, err := scanStudents(rows)
	if err != nil {
		return nil, err
	}

	if err := r.loadCourses(c, students); err != nil {
		return nil, err
	}
	return students, nil
}

func (r studentRepositorySQL) GetByID(c context.Context, id uint) (*model.Student, error) {
	var student model.Student

	err := r.db.QueryRowContext(c, `
		SELECT id, name, email
		FROM students
		WHERE id = $1
	`, id).Scan(&student.Id, &student.Name, &student.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get student %d: %w", id, err)
	}

	students := []model.Student{student}
	if err := r.loadCourses(c, students); err != nil {
		return nil, err
	}
	return &students[0], nil
}

func (r studentRepositorySQL) Create(c context.Context, student *model.Student) error {
	err := r.db.QueryRowContext(c, `
		INSERT INTO students (name, email)
		VALUES ($1, $2)
		RETURNING id
	`, student.Name, student.Email).Scan(&student.Id)

	if isUniqueViolation(err) {
		return fmt.Errorf("%w: student already exists", model.ErrConflict)
	}
	if err != nil {
		return fmt.Errorf("create student: %w", err)
	}
	return nil
}

func (r studentRepositorySQL) Update(c context.Context, student *model.Student) error {
	result, err := r.db.ExecContext(c, `
		UPDATE students
		SET name = $1, email = $2
		WHERE id = $3
	`, student.Name, student.Email, student.Id)

	if isUniqueViolation(err) {
		return fmt.Errorf("%w: student already exists", model.ErrConflict)
	}

	if err != nil {
		return fmt.Errorf("update student %d: %w", student.Id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update student %d: %w", student.Id, err)
	}
	if affected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r studentRepositorySQL) Delete(c context.Context, id uint) error {
	result, err := r.db.ExecContext(c, `
		DELETE FROM students WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete student %d: %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete student %d: %w", id, err)
	}
	if affected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r studentRepositorySQL) ListUnenrolled(c context.Context) ([]model.Student, error) {
	rows, err := r.db.QueryContext(c, `
		SELECT s.id, s.name, s.email
		FROM students s
		WHERE NOT EXISTS (
			SELECT 1
			FROM course_students cs
			WHERE cs.student_id = s.id
		)
		ORDER BY s.id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list unenrolled students: %w", err)
	}
	defer rows.Close()

	return scanStudents(rows)
}

func (r studentRepositorySQL) loadCourses(c context.Context, students []model.Student) error {
	if len(students) == 0 {
		return nil
	}

	coursesByStudent := make(map[uint][]model.Course, len(students))
	args := make([]any, len(students))
	placeholders := make([]string, len(students))

	for i, student := range students {
		coursesByStudent[student.Id] = []model.Course{}
		args[i] = student.Id
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT
			cs.student_id,
			c.id, c.title, c.description, c.duration_in_hours,
			c.start_date, c.instructor
		FROM course_students cs
		JOIN courses c on c.id = cs.course_id
		WHERE cs.student_id IN (%s)
		ORDER BY cs.student_id, c.id
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(c, query, args...)
	if err != nil {
		return fmt.Errorf("load enrolled courses: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var studentID uint
		var course model.Course

		if err := rows.Scan(
			&studentID,
			&course.Id,
			&course.Title,
			&course.Description,
			&course.DurationInHours,
			&course.StartDate,
			&course.Instructor,
		); err != nil {
			return fmt.Errorf("scan enrolled course: %w", err)
		}

		coursesByStudent[studentID] = append(coursesByStudent[studentID], course)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate enrolled course: %w", err)
	}

	for i := range students {
		students[i].EnrolledCourses = coursesByStudent[students[i].Id]
	}

	return nil
}

func scanStudents(rows *sql.Rows) ([]model.Student, error) {
	students := make([]model.Student, 0)

	for rows.Next() {
		var student model.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email); err != nil {
			return nil, fmt.Errorf("scan students: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate students: %w", err)
	}
	return students, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}