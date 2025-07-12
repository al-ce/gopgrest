package service

import (
	"time"

	"ftrack/models"
	"ftrack/repository"
)

// Service handles business logic with retrieved repository data
type Service struct {
	repo repository.Repository
}

// NewService returns a new Service struct
func NewService(r repository.Repository) Service {
	return Service{
		repo: r,
	}
}

// CreateSet inserts an exercise set in the exercise_sets table
func (s *Service) CreateSet(setData models.ExerciseSet) error {
	// Set time performed if not set
	if setData.PerformedAt.IsZero() {
		setData.PerformedAt = time.Now()
	}
	// Round to nearest second
	setData.PerformedAt = setData.PerformedAt.Round(time.Second)

	return s.repo.CreateSet(&setData)
}
