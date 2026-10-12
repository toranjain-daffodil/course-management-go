package model

import "time"

type Course struct {
	Id              int64 `json:"id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	DurationInHours int `json:"durationInHours" binding:"required"`
	StartDate       time.Time `json:"startDate" binding:"required"`
	Instructor 		string `json:"instructor" binding:"required"`
}