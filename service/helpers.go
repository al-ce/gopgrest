package service

import (
	"database/sql"
	"fmt"
	"reflect"

	"ftrack/repository"
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

	// Make a slice of zero values for this table and a slice of pointers to
	// those zero values
	rowValues, rowPtrs := makeScanDestination(table)

	scannedRowMapSlice := []types.RowDataMap{}
	for rows.Next() {
		// Scan column values into pointer slice
		err := rows.Scan(rowPtrs...)
		if err != nil {
			return []types.RowDataMap{}, err
		}
		// Create row data map of column names and column values
		scannedRowMap := makeScannedRowMap(table, rowValues)
		scannedRowMapSlice = append(scannedRowMapSlice, scannedRowMap)
	}

	return scannedRowMapSlice, nil
}

// scanSingleRow scans a single row from a query into a map
func (s *Service) scanSingleRow(tableName string, row *sql.Row) (types.RowDataMap, error) {
	// Get Table from Repository
	table, err := s.repo.GetTable(tableName)
	if err != nil {
		return types.RowDataMap{}, err
	}

	// Make a slice of zero values for this table and a slice of pointers to
	// those zero values
	rowValues, rowPtrs := makeScanDestination(table)

	// Scan column values into pointer slice
	err = row.Scan(rowPtrs...)
	if err != nil {
		return types.RowDataMap{}, err
	}
	// Create row data map of column names and column values
	scannedRowMap := makeScannedRowMap(table, rowValues)
	return scannedRowMap, nil
}

// makeScanDestination create slices to hold zero values for a given table and
// a slice of pointers to those zero values
func makeScanDestination(table *repository.Table) ([]any, []any) {
	rowValues := make([]any, len(table.Columns))
	rowPtrs := make([]any, len(table.Columns))
	for i := range rowValues {
		rowValues[i] = reflect.Zero(table.Columns[i].Datatype.ScanType())
		rowPtrs[i] = &rowValues[i]
	}
	return rowValues, rowPtrs
}

// makeScannedRowMap fills a RowDataMap with values from a scanned row
func makeScannedRowMap(table *repository.Table, rowValues []any) types.RowDataMap {
	scannedRowMap := types.RowDataMap{}
	for i := range len(table.Columns) {
		col := table.Columns[i].Name
		scannedRowMap[col] = rowValues[i]

	}
	return scannedRowMap
}
