package main

import (
	"assignment/internal/config"
	"assignment/internal/handler"
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

	gormDB, sqlDB, err := config.ConnectDatabase(dsn)
	if err != nil {
		log.Fatal(err)
	}
	courseRepository := repository.NewCourseRepository(gormDB, sqlDB)
	courseService := service.NewCourseService(courseRepository)
	courseHandler := handler.NewCourseHandler(courseService)
	router := router.SetupRouter(courseHandler)

	router.Run(":8080")
}