package service

import (
	"database/sql"
	"encoding/json"
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

// UpdateSet updates any number of valid ExerciseSet fields with separate calls
// to Repository.UpdateSetField
func (s *Service) UpdateSet(id string, updateData map[string]any) error {
	// Decode request body into a dummy ExerciseSet value to validate fields
	var dummyExerciseSet models.ExerciseSet
	b, _ := json.Marshal(updateData)
	err := json.Unmarshal(b, &dummyExerciseSet)
	if err != nil {
		return err
	}

	// Call Repository.UpdateSetField for each field to update
	for field, val := range updateData {
		if err := s.repo.UpdateSetField(id, field, val); err != nil {
			return err
		}
	}
	return nil
}

// DeleteSet removes a set from the exercise_sets table by id
func (s *Service) DeleteSet(id string) error {
	return s.repo.DeleteSet(id)
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
