package scheduler

import (
	"assignment/internal/service/interfaces"
	"context"
	"log"
	"time"
)

type SummaryScheduler struct {
	service interfaces.CourseService
	logger *log.Logger
}

func NewSummaryScheduler(service interfaces.CourseService, logger *log.Logger) *SummaryScheduler {
	return &SummaryScheduler{service: service, logger: logger}
}

func (s *SummaryScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <- ticker.C:
				summary, err := s.service.Summary(ctx)
				if err != nil {
					s.logger.Printf("summary job failed: %v", err)
					continue
				}

				s.logger.Printf("summary: courses=%d students=%d average_duration_hours=%.2f",
					summary.TotalCourses,
					summary.TotalStudents,
					summary.AverageDurationHours)
			case <-ctx.Done():
				s.logger.Println("summary scheduler stopped")
				return
			}
		}
	}()
}