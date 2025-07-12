package service

import (
	"database/sql"
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

// ListSets retrieves the list of all sets from the exercise_sets table
func (s *Service) ListSets() ([]models.ExerciseSet, error) {
	rows, err := s.repo.ListSets()
	if err != nil {
		return []models.ExerciseSet{}, err
	}
	defer rows.Close()

	// Scan rows into struct slice
	sets, err := scanExerciseSetRows(rows)
	if err != nil {
		return []models.ExerciseSet{}, err
	}
	return sets, nil
}

// scanExerciseSetRows scans sql rows into a slice of ExerciseSet structs
func scanExerciseSetRows(rows *sql.Rows) ([]models.ExerciseSet, error) {
	sets := []models.ExerciseSet{}
	for rows.Next() {
		set := &models.ExerciseSet{}
		err := rows.Scan(
			&set.ID,
			&set.Name,
			&set.PerformedAt,
			&set.Weight,
			&set.Unit,
			&set.Reps,
			&set.SetCount,
			&set.Notes,
			&set.SplitDay,
			&set.Program,
			&set.Tags,
		)
		if err != nil {
			return sets, err
		}
		sets = append(sets, *set)
	}
	return sets, nil
}
