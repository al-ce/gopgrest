package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gopgrest/apperrors"
	"gopgrest/rsql"
	"gopgrest/types"
)

// ListRowsByRSQL gets rows from a table with optional filter params
func (r *Repository) ListRowsByRSQL(tableName string, rsql *rsql.Query) (*sql.Rows, error) {
	// Build list of columns to select
	cols := buildColumnsToReturn(rsql)

	// Build list query with optional WHERE conditional filters
	conditional, values, err := buildWhereConditions(rsql.Filters)
	if err != nil {
		return nil, err
	}

	// Build list of optional JOIN relations
	joins := buildJoinRelations(rsql)

	listStmt := fmt.Sprintf("SELECT %s FROM %s %s %s", cols, tableName, joins, conditional)
	log.Printf("Exec query\n\t%s\nValues: %v\n", listStmt, values)

	// Execute list query
	rows, err := r.DB.Query(listStmt, values...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetRowByID gets a row from a table by id
func (r *Repository) GetRowByID(tableName string, id int64) (*sql.Rows, error) {
	log.Printf(
		"Exec query\n\tSELECT * FROM %s WHERE id = %d",
		tableName, id,
	)

	return r.DB.Query(
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", tableName),
		id,
	)
}

// InsertRow inserts a new row into a specified table
func (r *Repository) InsertRow(tableName string, newRow *types.RowData) (result ExecResult) {
	// Create cols/values/placeholders slices in consistent order
	var cols []string
	var values []any
	var valuesLog []string
	var placeholders []string
	var i int
	for k, v := range *newRow {
		cols = append(cols, k)
		values = append(values, v)
		valuesLog = append(valuesLog, fmt.Sprintf("%v", v))
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		i++
	}

	// Build create query
	createStmnt := fmt.Sprintf("INSERT INTO %s (", tableName) +
		strings.Join(cols, ", ") +
		") values (" +
		strings.Join(placeholders, ",") +
		") RETURNING id"

	log.Printf("Exec query\n\t%s\nValues: %v\n", createStmnt, values)

	// Execute insert query
	row := r.DB.QueryRow(createStmnt, values...)
	result.Error = row.Scan(&result.ID)
	return
}

// UpdateRowByID update columns in a table row by id
func (r *Repository) UpdateRowByID(tableName string, id int64, updatedRow *types.RowData) error {
	// Build update query
	updateStmnt := fmt.Sprintf("UPDATE %s SET ", tableName)

	// Create cols/values/placeholders slices in consistent order
	var values []any
	var i int
	for k, v := range *updatedRow {
		updateStmnt += fmt.Sprintf("%s = $%d, ", k, i+1)
		values = append(values, v)
		i++
	}
	// Strip final comma separator
	updateStmnt = updateStmnt[:len(updateStmnt)-2]
	updateStmnt += fmt.Sprintf(" WHERE id = %d", id)

	log.Printf("Exec query\n\t%s\nValues: %v", updateStmnt, values)

	// Execute update query
	result, err := r.DB.Exec(updateStmnt, values...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.NewUpdateInvalidIDErr(tableName, id)
	}

	return nil
}

// DeleteRowByID removes a row from a table by id
func (r *Repository) DeleteRowByID(tableName string, id int64) (int64, error) {
	deleteStmt := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	return r.execDeleteQuery(deleteStmt, []any{id})
}

// DeleteRowsByRSQL removes any rows matching the Filter in the Query
func (r *Repository) DeleteRowsByRSQL(tableName string, filters []rsql.Filter) (int64, error) {
	// Do not exec delete with empty query
	if len(filters) == 0 {
		return -1, apperrors.DeleteWithNoFilters
	}
	conditional, values, err := buildWhereConditions(filters)
	if err != nil {
		return -1, err
	}
	deleteStmt := fmt.Sprintf("DELETE FROM %s %s", tableName, conditional)
	return r.execDeleteQuery(deleteStmt, values)
}

func (r *Repository) execDeleteQuery(deleteStmt string, values []any) (int64, error) {
	log.Printf("Exec query\n\t%s\nValues: %v\n", deleteStmt, values)
	// Execute delete query
	result, err := r.DB.Exec(deleteStmt, values...)
	if err != nil {
		return -1, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowsAffected, nil
}
