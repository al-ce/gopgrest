package service

import (
	"database/sql"
	"fmt"
	"reflect"

	"ftrack/types"
)

// verifyColumns checks that all keys in a slice of cols, representing columns
// in a database table, actually exist in that table
func (s *Service) verifyColumns(tableName string, cols []string) error {
	table, err := s.repo.GetTable(tableName)
	if err != nil {
		return err
	}
	for _, col := range cols {
		if _, ok := table.ColumnMap[col]; !ok {
			return fmt.Errorf("Column '%s' does not exist in table %s", col, tableName)
		}
	}
	return nil
}

// scanRows scans rows from a query into a map
func (s *Service) scanRows(tableName string, rows *sql.Rows) ([]types.RowDataMap, error) {
	// Get Table from Repository
	table, err := s.repo.GetTable(tableName)
	if err != nil {
		return []types.RowDataMap{}, err
	}

	// Create slice to hold pointers of type/size equivalent to column type
	rowValues := make([]any, len(table.Columns))
	rowPtrs := make([]any, len(table.Columns))
	for i := range rowValues {
		rowValues[i] = reflect.Zero(table.Columns[i].Datatype.ScanType())
		rowPtrs[i] = &rowValues[i]
	}

	listQueryResults := []types.RowDataMap{}

	// Scan column values into pointer slice
	for rows.Next() {
		err := rows.Scan(rowPtrs...)
		if err != nil {
			return []types.RowDataMap{}, err
		}
		// Create map of column names and column values
		scannedRow := types.RowDataMap{}
		for i := range len(table.Columns) {
			col := table.Columns[i].Name
			scannedRow[col] = rowValues[i]

		}
		listQueryResults = append(listQueryResults, scannedRow)
	}

	return listQueryResults, nil
}
