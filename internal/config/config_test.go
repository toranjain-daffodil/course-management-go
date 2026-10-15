package config

import (
	"assignment/internal/model"
	"os"
	"testing"
)

func setupConfigTest(t *testing.T) {
	t.Helper()
	dir := os.TempDir()
	t.Chdir(dir)

	if err := os.WriteFile(".env", []byte(""), 0600); err != nil {
		t.Fatalf("failed to create .env: %v", err)
	}

	t.Setenv("POSTGRES_DSN", "test-dsn")
	t.Setenv("MAX_STUDENTS_PER_COURSE", "30")
	t.Setenv("MAX_COURSES_PER_STUDENT", "5")
	t.Setenv("DEFAULT_LIMIT", "10")
	t.Setenv("MAX_LIMIT", "100")
}

func TestNew_Success(t *testing.T) {
	setupConfigTest(t)
	cfg, err := New()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.PostgresDSN != "test-dsn" {
		t.Errorf("expected DSN %q, got %q", "test-dsn", cfg.PostgresDSN)
	}
	
	if cfg.Enrollment.MaxStudentsPerCourse != 30 {
		t.Errorf("expected MaxStudentsPerCourse 30, got %d",
			cfg.Enrollment.MaxStudentsPerCourse)
	}

	if cfg.Enrollment.MaxCoursesPerStudent != 5 {
		t.Errorf("expected MaxCoursesPerStudent 5, got %d",
			cfg.Enrollment.MaxCoursesPerStudent)
	}

	if cfg.Pagination.DefaultLimit != 10 {
		t.Errorf("expected DefaultLimit 10, got %d",
			cfg.Pagination.DefaultLimit)
	}

	if cfg.Pagination.MaxLimit != 100 {
		t.Errorf("expected MaxLimit 100, got %d",
			cfg.Pagination.MaxLimit)
	}
}

func TestNew_MissingDSN(t *testing.T) {
	setupConfigTest(t)
	t.Setenv("POSTGRES_DSN", "")

	cfg, err := New()

	if cfg != nil {
		t.Errorf("expected nil config, got %+v", cfg)
	}

	if err != model.ErrMissingConfig {
		t.Errorf("expected ErrMissingConfig, got %v", err)
	}
}

func TestNew_InvalidNumericConfig(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "invalid max students",
			key:   "MAX_STUDENTS_PER_COURSE",
			value: "abc",
		},
		{
			name:  "invalid max courses",
			key:   "MAX_COURSES_PER_STUDENT",
			value: "abc",
		},
		{
			name:  "invalid default limit",
			key:   "DEFAULT_LIMIT",
			value: "abc",
		},
		{
			name:  "invalid max limit",
			key:   "MAX_LIMIT",
			value: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupConfigTest(t)
			t.Setenv(tt.key, tt.value)

			cfg, err := New()

			if err == nil {
				t.Fatal("expected an error for invalid numeric configuration")
			}

			if cfg != nil {
				t.Errorf("expected nil config, got %+v", cfg)
			}
		})
	}
}