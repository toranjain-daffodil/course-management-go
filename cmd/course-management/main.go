package main

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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	
	dsn := os.Getenv("POSTGRES_DSN")

	gormDB, _, err := config.ConnectDatabase(dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := gormDB.AutoMigrate(&model.Course{}, &model.Student{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	courseRepository := repository.NewCourseRepositoryGORM(gormDB)
	studentRepository := repository.NewStudentRepositoryGORM(gormDB)

	courseService := service.NewCourseService(courseRepository)
	studentService := service.NewStudentService(studentRepository)

	courseHandler := handler.NewCourseHandler(courseService)
	studentHandler := handler.NewStudentHandler(studentService)

	router := router.SetupRouter(courseHandler, studentHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	summaryScheduler := scheduler.NewSummaryScheduler(courseService, log.Default())
	summaryScheduler.Start(ctx)

	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5 * time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}
}