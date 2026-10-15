package app

import (
	"assignment/internal/config"
	"assignment/internal/handler"
	"assignment/internal/model"
	"assignment/internal/repository"
	"assignment/internal/router"
	"assignment/internal/scheduler"
	"assignment/internal/service"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type App struct {
	server    *http.Server
	scheduler *scheduler.SummaryScheduler
	sqlDB     *sql.DB
}

func NewApp() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	gormDB, sqlDB, err := config.ConnectDatabase(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	if err := gormDB.AutoMigrate(&model.Course{}, &model.Student{}); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	courseRepository := repository.NewCourseRepositorySQL(sqlDB, cfg.Enrollment)
	studentRepository := repository.NewStudentRepositorySQL(sqlDB)

	courseService := service.NewCourseService(courseRepository, cfg.Pagination)
	studentService := service.NewStudentService(studentRepository, cfg.Pagination)

	courseHandler := handler.NewCourseHandler(courseService)
	studentHandler := handler.NewStudentHandler(studentService)

	router := router.SetupRouter(courseHandler, studentHandler)

	summaryScheduler := scheduler.NewSummaryScheduler(courseService, log.Default())

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	app := &App{server: server, scheduler: summaryScheduler, sqlDB: sqlDB}
	return app, nil
}

func (app *App) Run(ctx context.Context) (runErr error) {
	runCtx, cancel := context.WithCancel(ctx)
	schedulerDone := app.scheduler.Start(runCtx)

	defer func() {
		cancel()

		<-schedulerDone

		if err := app.sqlDB.Close(); err != nil && runErr == nil {
			runErr = fmt.Errorf("close database pool: %w", err)
		}
	}()

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := app.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		return nil

	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", err)
		}

		return nil
	}
}
