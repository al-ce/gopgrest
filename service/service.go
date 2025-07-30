package service

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"ftrack/repository"
	"ftrack/types"
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
func (s *Service) InsertRow(newRow *types.RowData, tableName string) (int64, error) {
	// Each column in the insert data must exist in the table
	cols := slices.Collect(maps.Keys(*newRow))
	if err := s.verifyColumns(tableName, cols); err != nil {
		return -1, err
	}

	result := s.repo.InsertRow(tableName, newRow)
	return result.ID, result.Error
}

// PickRow gets a row from a table by id
func (s *Service) PickRow(tableName, id string) (types.RowData, error) {
	row := s.repo.GetRowByID(tableName, id)
	gotId, rowData, err := s.scanSingleRow(tableName, row)
	if fmt.Sprintf("%v", gotId) != id {
		return types.RowData{}, fmt.Errorf("PickRow got id %v, requested %s", gotId, id)
	}
	return rowData, err
}

// ListRows gets rows from a table with optional filter params
func (s *Service) ListRows(tableName string, qf types.QueryFilter) (types.RowDataIdMap, error) {
	// Each column in the query params must exist in the table
	cols := slices.Collect(maps.Keys(qf))
	if err := s.verifyColumns(tableName, cols); err != nil {
		return types.RowDataIdMap{}, err
	}

	rows, err := s.repo.ListRows(tableName, qf)
	if err != nil {
		return types.RowDataIdMap{}, err
	}
	defer rows.Close()

	// Scan rows into struct slice
	listQueryResults, err := s.scanRows(tableName, rows)
	if err != nil {
		return types.RowDataIdMap{}, err
	}
	return listQueryResults, nil
}

// UpdateRow updates any number of valid fields with separate calls to
// Repository.UpdateRowCol
func (s *Service) UpdateRow(tableName, id string, updateData types.RowData) error {
	// Each column in the update data must exist in the table
	cols := slices.Collect(maps.Keys(updateData))
	if err := s.verifyColumns(tableName, cols); err != nil {
		return err
	}

	// Decode request body into a dummy row value to validate fields
	var dummyRow types.RowData
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
