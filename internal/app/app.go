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
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type App struct {
	server    *http.Server
	scheduler *scheduler.SummaryScheduler
}

func NewApp() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	gormDB, _, err := config.ConnectDatabase(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	if err := gormDB.AutoMigrate(&model.Course{}, &model.Student{}); err != nil {
		return nil, err
	}

	courseRepository := repository.NewCourseRepositoryGORM(gormDB, cfg.Enrollment)
	studentRepository := repository.NewStudentRepositoryGORM(gormDB)

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

	app := &App{server: server, scheduler: summaryScheduler}
	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	app.scheduler.Start(ctx)
	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %v", err)
		}
	case err := <-serverErrors:
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %v", err)
		}
	}
	return nil
}
