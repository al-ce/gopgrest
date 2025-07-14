package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

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
func (s *Service) InsertRow(newRow *map[string]any, table string) error {
	if !s.repo.TableExists(table) {
		return fmt.Errorf("table %s does not exist", table)
	}
	return s.repo.InsertRow(table, newRow)
}

// ListRows gets rows from a table with optional filter params
func (s *Service) ListRows(table string, params map[string][]string) ([]map[string]any, error) {
	rows, err := s.repo.ListRows(table, params)
	if err != nil {
		return []map[string]any{}, err
	}
	defer rows.Close()

	// Scan rows into struct slice
	listQueryResults, err := scanRows(rows)
	if err != nil {
		return []map[string]any{}, err
	}
	return listQueryResults, nil
}

// UpdateRow updates any number of valid fields with separate calls to
// Repository.UpdateRowCol
func (s *Service) UpdateRow(table, id string, updateData map[string]any) error {
	// Decode request body into a dummy row value to validate fields
	var dummyRow map[string]any
	b, _ := json.Marshal(updateData)
	err := json.Unmarshal(b, &dummyRow)
	if err != nil {
		return err
	}

	// Update individual columns in the row
	for field, val := range updateData {
		if err := s.repo.UpdateRowCol(table, id, field, val); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRow removes a row form the table by id
func (s *Service) DeleteRow(table, id string) error {
	return s.repo.DeleteRow(table, id)
}

// scanRows scans rows from a query into a map
func scanRows(rows *sql.Rows) ([]map[string]any, error) {
	// Get column names
	cols, err := rows.Columns()
	if err != nil {
		return []map[string]any{}, err
	}

	// Get column types
	coltypes, err := rows.ColumnTypes()
	if err != nil {
		return []map[string]any{}, err
	}

	// Create slice to hold pointers of type/size equivalent to column type
	rowValues := make([]any, len(cols))
	rowPtrs := make([]any, len(cols))
	for i := range rowValues {
		rowValues[i] = reflect.Zero(coltypes[i].ScanType())
		rowPtrs[i] = &rowValues[i]
	}

	listQueryResults := []map[string]any{}

	// Scan column values into pointer slice
	for rows.Next() {
		err := rows.Scan(rowPtrs...)
		if err != nil {
			return []map[string]any{}, err
		}
		// Create map of column names and column values
		scannedRow := map[string]any{}
		for i := range len(cols) {
			col := cols[i]
			scannedRow[col] = rowValues[i]

		}
		listQueryResults = append(listQueryResults, scannedRow)
	}

	return listQueryResults, nil
}
