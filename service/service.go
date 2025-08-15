package service

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"
	"strconv"

	"gopgrest/repository"
	"gopgrest/rsql"
	"gopgrest/types"
)

// Service handles business logic with retrieved repository data
type Service struct {
	Repo repository.Repository
}

// NewService returns a new Service struct
func NewService(r repository.Repository) Service {
	return Service{
		Repo: r,
	}
}

// InsertRow inserts a new row in a specified table
func (s *Service) InsertRow(newRow *types.RowData, tableName string) (int64, error) {
	table, err := s.Repo.GetTable(tableName)
	if err != nil {
		return -1, err
	}

	// Each column in the insert data must exist in the table
	cols := slices.Collect(maps.Keys(*newRow))
	if err := verifyColumns(table, cols); err != nil {
		return -1, err
	}

	result := s.Repo.InsertRow(tableName, newRow)
	return result.ID, result.Error
}

// PickRow gets a row from a table by id
func (s *Service) PickRow(tableName, id string) (types.RowData, error) {
	// Get table info for verification
	_, err := s.Repo.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	// Parse ID
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	// Get Row from database, expect 1 row
	rows, err := s.Repo.GetRowByID(tableName, idInt)
	if err != nil {
		return nil, err
	}
	rowData, err := ScanRows(rows)
	if err != nil {
		return nil, err
	}
	if len(rowData) != 1 {
		return nil, fmt.Errorf("Expected 1 result for pick, got %d", len(rowData))
	}
	return rowData[0], err
}

// ListRows gets rows from a table with optional filter params
func (s *Service) ListRows(tableName string, url string) ([]types.RowData, error) {
	// Get table info for verification
	_, err := s.Repo.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	// Parse RSQL
	rsql, err := rsql.NewRSQLQuery(url)
	if err != nil {
		return nil, err
	}
	// Validate RSQL
	if err := s.validateRSQLQuery(rsql); err != nil {
		return nil, err
	}

	// Query db
	rows, err := s.Repo.ListRows(tableName, rsql)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %s\n", err)
		}
	}()

	// Scan rows into struct slice
	listQueryResults, err := ScanRows(rows)
	if err != nil {
		return nil, err
	}
	return listQueryResults, nil
}

// UpdateRow updates any number of valid columns with separate calls to
// Repository.UpdateRowCol
func (s *Service) UpdateRow(tableName, id string, updateData *types.RowData) (types.RowData, error) {
	table, err := s.Repo.GetTable(tableName)
	if err != nil {
		return types.RowData{}, err
	}

	// Convert id to int
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return types.RowData{}, err
	}

	// Each column in the update data must exist in the table
	cols := slices.Collect(maps.Keys(*updateData))
	if err := verifyColumns(table, cols); err != nil {
		return types.RowData{}, err
	}

	// Decode request body into a dummy row value to validate column names
	var dummyRow types.RowData
	b, _ := json.Marshal(updateData)
	err = json.Unmarshal(b, &dummyRow)
	if err != nil {
		return types.RowData{}, err
	}

	// Update row
	err = s.Repo.UpdateRow(tableName, idInt, updateData)
	if err != nil {
		return types.RowData{}, err
	}
	return s.PickRow(tableName, id)
}

// DeleteRow removes a row from the table by id
func (s *Service) DeleteRow(tableName, id string) (int64, error) {
	// Get table info for verification
	_, err := s.Repo.GetTable(tableName)
	if err != nil {
		return -1, err
	}

	// Convert id to int
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return -1, err
	}

	return s.Repo.DeleteRow(tableName, idInt)
}
