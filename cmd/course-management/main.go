package main

import (
	"assignment/internal/config"
	"assignment/internal/handler"
	"assignment/internal/model"
	"assignment/internal/repository"
	"assignment/internal/router"
	"assignment/internal/service"
	"log"
	"os"

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

	router.Run(":8080")
}