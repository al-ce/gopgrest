package service

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"
	"strconv"

	"gopgrest/apperrors"
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

// GetRowByID gets a row from a table by id
func (s *Service) GetRowByID(tableName, idAsStr string) (types.RowData, error) {
	// Get table info for verification
	_, err := s.Repo.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	// Parse ID
	idAsInt, err := strconv.ParseInt(idAsStr, 10, 64)
	if err != nil {
		return nil, err
	}
	// Get Row from database, expect 1 row
	rows, err := s.Repo.GetRowByID(tableName, idAsInt)
	if err != nil {
		return nil, err
	}
	rowData, err := ScanRows(rows)
	if err != nil {
		return nil, err
	}
	if len(rowData) != 1 {
		return types.RowData{}, fmt.Errorf("%s (id: %s)", apperrors.GetByIdNotUnique, idAsStr)
	}
	return rowData[0], err
}

// GetRowsByRSQL gets rows from a table with optional filter params
func (s *Service) GetRowsByRSQL(tableName string, url string) ([]types.RowData, error) {
	// Get table info for verification
	_, err := s.Repo.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	// Parse RSQL
	query, err := rsql.NewRSQLQuery(url)
	if err != nil {
		return nil, err
	}
	// Validate RSQL
	if err := s.validateRSQLQuery(query); err != nil {
		return nil, err
	}

	// Query db
	rows, err := s.Repo.GetRowsByRSQL(tableName, query)
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

// InsertRows inserts new rows in a specified table If multiple rows are
// inserted, they must each have the same columns and value types
func (s *Service) InsertRows(newRows []types.RowData, tableName string) ([]int64, error) {
	var ids []int64
	if len(newRows) == 0 {
		return ids, apperrors.InsertWithNoRows
	}
	table, err := s.Repo.GetTable(tableName)
	if err != nil {
		return ids, err
	}

	// Each column in the insert data must exist in the table
	cols := slices.Collect(maps.Keys(newRows[0]))
	badCol, err := verifyColumns(table, cols)
	if err != nil {
		return ids, fmt.Errorf("%w\n(%s:%s) ", err, table.Name, badCol)
	}

	// If there's only one row to insert, skip the remaining consistency checks
	if len(newRows) == 1 {
		return s.Repo.InsertRows(tableName, newRows)
	}

	// Compare cols and value types of other rows against first row
	slices.Sort(cols)
	valTypes := make(map[string]any, len(cols))
	for _, col := range cols {
		valTypes[col] = fmt.Sprintf("%T", newRows[0][col])
	}

	// Each row must have matching columns and value types
	for _, row := range newRows {
		// Compare cols
		thisCols := slices.Collect(maps.Keys(row))
		slices.Sort(thisCols)
		if slices.Compare(cols, thisCols) != 0 {
			return ids, fmt.Errorf("%w\n%v %v", apperrors.InsertColsDoNotMatch, cols, thisCols)
		}
		// Compare value types
		for k, v := range row {
			if valTypes[k] != fmt.Sprintf("%T", v) {
				return ids, fmt.Errorf(
					"%w\ncol: %s %v %v",
					apperrors.InsertValTypesDoNotMatch,
					k,
					valTypes[k],
					v,
				)
			}
		}
	}

	return s.Repo.InsertRows(tableName, newRows)
}

// UpdateRowByID updates any number of valid columns with separate calls to
// Repository.UpdateRowCol
func (s *Service) UpdateRowByID(tableName, id string, updateData *types.RowData) (types.RowData, error) {
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
	badCol, err := verifyColumns(table, cols)
	if err != nil {
		return types.RowData{}, fmt.Errorf("%w (%s:%s) ", err, table.Name, badCol)
	}

	// Decode request body into a dummy row value to validate column names
	var dummyRow types.RowData
	b, _ := json.Marshal(updateData)
	err = json.Unmarshal(b, &dummyRow)
	if err != nil {
		return types.RowData{}, err
	}

	// Update row
	err = s.Repo.UpdateRowByID(tableName, idInt, updateData)
	if err != nil {
		return types.RowData{}, err
	}
	return s.GetRowByID(tableName, id)
}

// DeleteRowByID removes a row from the table by id
func (s *Service) DeleteRowByID(tableName, id string) (int64, error) {
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

	return s.Repo.DeleteRowByID(tableName, idInt)
}
