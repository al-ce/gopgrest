package service

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"ftrack/repository"
	"ftrack/types"
)

// verifyColumns checks that all keys in a slice of cols, representing columns
// in a database table, actually exist in that table
func (s *Service) verifyColumns(t *repository.Table, cols []string) error {
	for _, col := range cols {
		if _, ok := t.ColumnMap[col]; !ok {
			return fmt.Errorf("Column '%s' does not exist in table %s", col, t.Name)
		}
	}
	return nil
}

// scanRows scans rows from a query into a map
func (s *Service) scanRows(tableName string, rows *sql.Rows) (*types.RowDataIdMap, error) {
	// Get Table from Repository
	table, err := s.Repo.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	// Make a slice of zero values for this table and a slice of pointers to
	// those zero values
	rowValues, rowPtrs := makeScanDestination(table)

	scannedRowIdMap := make(types.RowDataIdMap)
	for rows.Next() {
		// Scan column values into pointer slice
		err := rows.Scan(rowPtrs...)
		if err != nil {
			return nil, err
		}
		// Create row data map of column names and column values
		id, scannedRow, err := makeScannedRowMap(table, rowValues)
		if err != nil {
			return nil, nil
		}
		scannedRowIdMap[id] = scannedRow
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows in table %s: %v", tableName, err)
		return nil, err
	}

	return &scannedRowIdMap, nil
}

// scanSingleRow scans a single row from a query into a map
func (s *Service) scanSingleRow(tableName string, row *sql.Row) (
	int64,
	types.RowData,
	error,
) {
	// Get Table from Repository
	table, err := s.Repo.GetTable(tableName)
	if err != nil {
		return -1, types.RowData{}, err
	}

	// Make a slice of zero values for this table and a slice of pointers to
	// those zero values
	rowValues, rowPtrs := makeScanDestination(table)

	// Scan column values into pointer slice
	err = row.Scan(rowPtrs...)
	if err != nil {
		return -1, types.RowData{}, err
	}
	// Create row data map of column names and column values
	id, scannedRowMap, err := makeScannedRowMap(table, rowValues)
	if err != nil {
		return -1, types.RowData{}, err
	}
	return id, scannedRowMap, nil
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
func makeScannedRowMap(table *repository.Table, rowValues []any) (
	int64,
	types.RowData,
	error,
) {
	var id int64
	scannedRow := make(types.RowData)
	for i := range len(table.Columns) {
		col := table.Columns[i].Name
		val := rowValues[i]
		if col == "id" && reflect.TypeOf(val) == reflect.TypeOf(int64(0)) {
			idCast, isInt64 := val.(int64)
			if !isInt64 {
				return -1, types.RowData{}, fmt.Errorf("makeScannedRowMap could not convert ")
			} else {
				id = idCast
			}
		}
		scannedRow[col] = val

	}
	return id, scannedRow, nil
}
