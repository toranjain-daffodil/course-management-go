package repository

import (
	"assignment/internal/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type courseRepositoryGORM struct {
	db *gorm.DB
	enroll model.EnrollmentConfig
}

func NewCourseRepositoryGORM(
	gormDB *gorm.DB,
	enroll model.EnrollmentConfig,
) *courseRepositoryGORM {
	return &courseRepositoryGORM{
		db: gormDB,
		enroll: enroll,
	}
}

func (repo courseRepositoryGORM) List(c context.Context, filter model.CourseFilter) ([]model.Course, error) {
	q := repo.db.WithContext(c).Model(&model.Course{})

	if filter.Search != "" {
		q = q.Where("title ILIKE ?", "%"+filter.Search+"%")
	}

	var courses []model.Course
	err := q.Preload("Students").Order("id ASC").Limit(filter.Limit).Offset(filter.Offset).Find(&courses).Error
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	return courses, nil
}

func (repo courseRepositoryGORM) GetByID(c context.Context, id uint) (*model.Course, error) {
	var course model.Course
	err := repo.db.WithContext(c).Preload("Students").First(&course, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get course %d: %w", id, err)
	}
	return &course, nil
}

func (repo courseRepositoryGORM) Create(c context.Context, course *model.Course) error {
	err := repo.db.WithContext(c).Create(course).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: course title already exists", model.ErrConflict)
	}
	if err != nil {
		return fmt.Errorf("create course: %w", err)
	}
	return nil
}

func (repo courseRepositoryGORM) Update(c context.Context, course *model.Course) error {
	res := repo.
		db.
		WithContext(c).
		Model(course).
		Select(
			"Title",
			"Description",
			"DurationInHours",
			"StartDate",
			"Instructor",
		).
		Updates(course)
	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: course title already exists", model.ErrConflict)
	}
	if res.Error != nil {
		return fmt.Errorf("update course %d: %w", course.Id, res.Error)
	}
	if res.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (repo courseRepositoryGORM) Delete(c context.Context, id uint) error {
	res := repo.db.WithContext(c).Delete(&model.Course{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete course %d: %w", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (repo courseRepositoryGORM) Enroll(c context.Context, courseID uint, studentID uint) error {
	return repo.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var course model.Course
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&course, courseID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model.ErrNotFound
			}
			return fmt.Errorf("get course for enrollment: %w", err)
		}

		var student model.Student
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&student, studentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model.ErrNotFound
			}
			return fmt.Errorf("get student for enrollment: %w", err)
		}

		var existing int64
		if err := tx.Table("course_students").Where("course_id = ? AND student_id = ?", courseID, studentID).Count(&existing).Error; err != nil {
			return fmt.Errorf("check enrollment: %w", err)
		}
		if existing > 0 {
			return fmt.Errorf("%w: student is already enrolled in this course", model.ErrConflict)
		}

		var courseEnrollmentCount int64
		if err := tx.Table("course_students").Where("course_id = ?", courseID).Count(&courseEnrollmentCount).Error; err != nil {
			return fmt.Errorf("count course enrollments: %w", err)
		}
		if courseEnrollmentCount >= int64(repo.enroll.MaxStudentsPerCourse) {
			return fmt.Errorf("%w: a course cannot have more than %d students", model.ErrValidation, repo.enroll.MaxStudentsPerCourse)
		}

		var studentEnrollmentCount int64
		if err := tx.Table("course_students").Where("student_id = ?", studentID).Count(&studentEnrollmentCount).Error; err != nil {
			return fmt.Errorf("count student enrollments: %w", err)
		}
		if studentEnrollmentCount >= int64(repo.enroll.MaxCoursesPerStudent) {
			return fmt.Errorf("%w: a student cannot enroll in more than %d courses", model.ErrValidation, repo.enroll.MaxCoursesPerStudent)
		}

		if err := tx.Model(&course).Association("Students").Append(&student); err != nil {
			return fmt.Errorf("enroll student: %w", err)
		}
		return nil
	})
}

func (repo courseRepositoryGORM) TopEnrolled(c context.Context, limit int) ([]model.Course, error) {
	var courses []model.Course
	res := repo.db.WithContext(c).
		Model(&model.Course{}).
		Select("courses.*").
		Joins("LEFT JOIN course_students ON course_students.course_id = courses.id").
		Group("courses.id").
		Order("COUNT(course_students.student_id) DESC, courses.id ASC").
		Limit(limit).
		Preload("Students").
		Find(&courses)
	if res.Error != nil {
		return nil, fmt.Errorf("list top enrolled courses: %w", res.Error)
	}
	return courses, nil
}

func (repo courseRepositoryGORM) Summary(c context.Context) (model.CourseSummary, error) {
	var summary model.CourseSummary

	if err := repo.db.WithContext(c).
				Model(&model.Course{}).
				Count(&summary.TotalCourses).Error; err != nil {
					return summary, err
				}

	if err := repo.db.WithContext(c).
				Model(&model.Student{}).
				Count(&summary.TotalStudents).Error; err != nil {
					return summary, err
				}

	if err := repo.db.WithContext(c).
				Model(&model.Course{}).
				Select("COALESCE(AVG(duration_in_hours), 0)").
				Scan(&summary.AverageDurationHours).Error; err != nil {
					return summary, err
				}
	return summary, nil
}