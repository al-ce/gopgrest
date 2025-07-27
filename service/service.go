package service

import (
	"encoding/json"
	"maps"
	"slices"

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

// InsertRow inserts a new row in a specified table
func (s *Service) InsertRow(newRow *map[string]any, tableName string) (int64, error) {
	// Each column in the insert data must exist in the table
	cols := slices.Collect(maps.Keys(*newRow))
	if err := s.verifyColumns(tableName, cols); err != nil {
		return 0, err
	}

	return s.repo.InsertRow(tableName, newRow)
}

// ListRows gets rows from a table with optional filter params
func (s *Service) ListRows(tableName string, params map[string][]string) ([]map[string]any, error) {
	// Each column in the query params must exist in the table
	cols := slices.Collect(maps.Keys(params))
	if err := s.verifyColumns(tableName, cols); err != nil {
		return []map[string]any{}, err
	}

	rows, err := s.repo.ListRows(tableName, params)
	if err != nil {
		return []map[string]any{}, err
	}
	defer rows.Close()

	// Scan rows into struct slice
	listQueryResults, err := s.scanRows(tableName, rows)
	if err != nil {
		return []map[string]any{}, err
	}
	return listQueryResults, nil
}

// UpdateRow updates any number of valid fields with separate calls to
// Repository.UpdateRowCol
func (s *Service) UpdateRow(tableName, id string, updateData map[string]any) error {
	// Each column in the update data must exist in the table
	cols := slices.Collect(maps.Keys(updateData))
	if err := s.verifyColumns(tableName, cols); err != nil {
		return err
	}

	// Decode request body into a dummy row value to validate fields
	var dummyRow map[string]any
	b, _ := json.Marshal(updateData)
	err := json.Unmarshal(b, &dummyRow)
	if err != nil {
		return err
	}

	// Update individual columns in the row
	for field, val := range updateData {
		if err := s.repo.UpdateRowCol(tableName, id, field, val); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRow removes a row from the table by id
func (s *Service) DeleteRow(tableName, id string) error {
	return s.repo.DeleteRow(tableName, id)
}
