package service

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"gopgrest/repository"
	"gopgrest/rsql"
	"gopgrest/types"
)

// verifyColumns checks that all keys in a slice of cols, representing columns
// in a database table, actually exist in that table
// TODO: replace all these refs with s.Repo.IsValidColumn
func verifyColumns(t *repository.Table, cols []string) error {
	for _, col := range cols {
		if _, ok := t.ColumnMap[col]; !ok {
			return fmt.Errorf("Column '%s' does not exist in table %s", col, t.Name)
		}
	}
	return nil
}

// scanRows scans rows from a query into a map
func (s *Service) scanRows(rows *sql.Rows) ([]types.RowData, error) {
	// Make arrays of pointers with sizes that match column type
	cols, _ := rows.Columns()
	rowValues, rowPtrs := makeScanDestination(rows, cols)

	scannedRows := []types.RowData{}
	for rows.Next() {
		// Scan column values into pointer slice
		err := rows.Scan(rowPtrs...)
		if err != nil {
			return nil, err
		}
		scannedRow := makeScannedRowMap(cols, rowValues)
		scannedRows = append(scannedRows, scannedRow)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating over rows: %v", err)
		return nil, err
	}

	return scannedRows, nil
}

// makeScanDestination create slices to hold zero values for a given query and
// a slice of pointers to those zero values
func makeScanDestination(rows *sql.Rows, cols []string) ([]any, []any) {
	ct, _ := rows.ColumnTypes()
	rowValues := make([]any, len(cols))
	rowPtrs := make([]any, len(cols))
	for i, v := range ct {
		scanType := v.ScanType()
		zeroVal := reflect.Zero(scanType)
		rowValues[i] = zeroVal
		rowPtrs[i] = &rowValues[i]
	}
	return rowValues, rowPtrs
}

// makeScannedRowMap fills a RowDataMap with values from a scanned row
func makeScannedRowMap(cols []string, rowValues []any) types.RowData {
	scannedRow := make(types.RowData)
	for i, col := range cols {
		val := rowValues[i]
		scannedRow[col] = val

	}
	return scannedRow
}

func (s *Service) validateRSQLQuery(query *rsql.Query) error {
	if query == nil {
		return nil
	}
	if err := s.validateRSQLFilters(query.Filters); err != nil {
		return err
	}
	if err := s.validateRSQLFields(query.Fields); err != nil {
		return err
	}
	if err := s.validateRSQLJoins(query.Joins); err != nil {
		return err
	}
	return nil
}

func (s *Service) validateRSQLFilters(filters []rsql.Filter) error {
	// Validate: each column in the query filter should be valid for its table
	for _, f := range filters {
		// Check if column is prefixed with a table, e.g. authors.forename
		prefixedCol := strings.Split(f.Column, ".")
		if len(prefixedCol) == 2 {
			tableName := prefixedCol[0]
			colName := prefixedCol[1]
			table, err := s.Repo.GetTable(tableName)
			if err != nil {
				return err
			}
			if !s.Repo.IsValidColumn(*table, colName) {
				return fmt.Errorf(
					"Invalid col name %s for table %s",
					colName,
					tableName,
				)
			}
		} else {
			// TODO: get tables from JOINs
			// We don't know which table this col could belong to, so check any
			// tables from the url (the table in the FROM clause in the URL
			// resource + any tables mentioned in the JOIN params of the URL)
		}
	}
	return nil
}

func (s *Service) validateRSQLFields(fields []rsql.Field) error {
	for _, f := range fields {

		// If the field has a qualifier, check the column against that table
		if f.Qualifier != "" {
			t, err := s.Repo.GetTable(f.Qualifier)
			if err != nil {
				return err
			}
			if !s.Repo.IsValidColumn(*t, f.Column) {
				return fmt.Errorf("field %s not found in table %s", f, f.Qualifier)
			}
			continue
		}

		// Otherwise, check all tables
		foundCol := false
		for _, t := range s.Repo.Tables {
			if s.Repo.IsValidColumn(t, f.Column) {
				foundCol = true
				break
			}
		}
		if !foundCol {
			return fmt.Errorf("field %s not found in any tables", f)
		}
	}
	return nil
}

func (s *Service) validateRSQLJoins(joins []rsql.JoinRelation) error {
	for _, j := range joins {
		// Check tables exist
		if _, err := s.Repo.GetTable(j.Table); err != nil {
			return err
		}
		leftTable, err := s.Repo.GetTable(j.LeftQualifier)
		if err != nil {
			return err
		}
		rightTable, err := s.Repo.GetTable(j.RightQualifier)
		if err != nil {
			return err
		}
		// Check table after JOIN keyword is in the qualified column names
		// (LeftTable/RightTable)
		if !slices.Contains([]string{j.LeftQualifier, j.RightQualifier}, j.Table) {
			return fmt.Errorf("Table in JOIN statement missing from relation: %v", j)
		}

		if !s.Repo.IsValidColumn(*leftTable, j.LeftCol) {
			return fmt.Errorf("col %s not found in table %s in join %v", j.LeftCol, leftTable, j)
		}
		if !s.Repo.IsValidColumn(*rightTable, j.RightCol) {
			return fmt.Errorf("col %s not found in table %s in join %v", j.LeftCol, leftTable, j)
		}
	}
	return nil
}
