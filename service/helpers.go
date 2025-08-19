package service

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"gopgrest/apperrors"
	"gopgrest/repository"
	"gopgrest/rsql"
	"gopgrest/types"
)

// verifyColumns checks that all keys in a slice of cols, representing columns
// in a database table, actually exist in that table
func verifyColumns(t *repository.Table, cols []string) (string, error) {
	for _, col := range cols {
		if _, ok := t.ColumnMap[col]; !ok {
			return col, apperrors.ColDoesNotExist
		}
	}
	return "", nil
}

// ScanRows scans rows from a query into a map
func ScanRows(rows *sql.Rows) ([]types.RowData, error) {
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

func (s *Service) validateRSQLQuery(query rsql.QueryParams) error {
	if err := s.validateRSQLTables(query.Tables); err != nil {
		return err
	}
	if err := s.ValidateRSQLConditions(query.Tables, query.Conditions); err != nil {
		return err
	}
	if err := s.validateRSQLColumns(query.Columns); err != nil {
		return err
	}
	if err := s.validateRSQLJoins(query.Joins); err != nil {
		return err
	}
	return nil
}

func (s *Service) ValidateRSQLConditions(tableNames []string, conditions []rsql.Condition) error {
	// Validate: each column in the WHERE clause should be valid for its table
	for _, f := range conditions {
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
			// Search for column in all tables referenced in query
			found := false
			for _, tableName := range tableNames {
				table, err := s.Repo.GetTable(tableName)
				if err != nil {
					return err
				}
				if s.Repo.IsValidColumn(*table, f.Column) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf(
					"Could not find col %s in any referenced tables",
					f.Column,
				)
			}
		}
	}
	return nil
}

func (s *Service) validateRSQLTables(tables []string) error {
	for _, t := range tables {
		_, err := s.Repo.GetTable(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) validateRSQLColumns(columns []rsql.Column) error {
	for _, f := range columns {

		// If the column has a qualifier, check the column against that table
		if f.Qualifier != "" {
			t, err := s.Repo.GetTable(f.Qualifier)
			if err != nil {
				return err
			}
			if !s.Repo.IsValidColumn(*t, f.Name) {
				return fmt.Errorf("column %s not found in table %s", f, f.Qualifier)
			}
			continue
		}

		// Otherwise, check all tables
		foundCol := false
		for _, t := range s.Repo.Tables {
			if s.Repo.IsValidColumn(t, f.Name) {
				foundCol = true
				break
			}
		}
		if !foundCol {
			return fmt.Errorf("column %s not found in any tables", f)
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

func (s *Service) newRSQLQuery(url string) (rsql.QueryParams, error) {
	// Parse RSQL
	query, err := rsql.NewRSQLQuery(url)
	if err != nil {
		return rsql.QueryParams{}, err
	}
	// Validate RSQL
	if err := s.validateRSQLQuery(query); err != nil {
		return rsql.QueryParams{}, err
	}
	return query, nil
}
