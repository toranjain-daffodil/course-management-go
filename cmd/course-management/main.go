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

	if err := gormDB.AutoMigrate(&model.Course{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	courseRepository := repository.NewCourseRepositoryGORM(gormDB)
	courseService := service.NewCourseService(courseRepository)
	courseHandler := handler.NewCourseHandler(courseService)
	router := router.SetupRouter(courseHandler)

	router.Run(":8080")
}