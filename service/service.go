package service

import (
	"database/sql"
	"encoding/json"
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
func (s *Service) InsertRow(newRow *map[string]any, tableName string) error {
	if _, err := s.repo.GetTable(tableName); err != nil {
		return err
	}
	// TODO: verify new row cols are in the table

	return s.repo.InsertRow(tableName, newRow)
}

// ListRows gets rows from a table with optional filter params
func (s *Service) ListRows(tableName string, params map[string][]string) ([]map[string]any, error) {
	// TODO: verify param cols are in the table

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
	// TODO: verify update cols are in the table

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

// scanRows scans rows from a query into a map
func (s *Service) scanRows(tableName string, rows *sql.Rows) ([]map[string]any, error) {
	// Get Table from Repository
	table, err := s.repo.GetTable(tableName)
	if err != nil {
		return []map[string]any{}, err
	}

	// Create slice to hold pointers of type/size equivalent to column type
	rowValues := make([]any, len(table.Columns))
	rowPtrs := make([]any, len(table.Columns))
	for i := range rowValues {
		rowValues[i] = reflect.Zero(table.Columns[i].Datatype.ScanType())
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
		for i := range len(table.Columns) {
			col := table.Columns[i].Name
			scannedRow[col] = rowValues[i]

		}
		listQueryResults = append(listQueryResults, scannedRow)
	}

	return listQueryResults, nil
}
