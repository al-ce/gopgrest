package repository

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"

	"gopgrest/apperrors"
	"gopgrest/rsql"
	"gopgrest/types"
)

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

// GetRowsByRSQL gets rows from a table with optional query params
func (r *Repository) GetRowsByRSQL(tableName string, query rsql.QueryParams) (*sql.Rows, error) {
	// Build list of columns to select
	cols := buildColumnsToReturn(query)

	// Build list query with optional WHERE conditional statements
	conditional, values, err := buildWhereConditions(query.Conditions, 0)
	if err != nil {
		return nil, err
	}

	// Build list of optional JOIN relations
	joins := buildJoinRelations(query)

	listStmt := fmt.Sprintf("SELECT %s FROM %s %s %s", cols, tableName, joins, conditional)
	log.Printf("Exec query\n\t%s\nValues: %v\n", listStmt, values)

	// Execute list query
	rows, err := r.DB.Query(listStmt, values...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// InsertRows inserts a new row into a specified table
func (r *Repository) InsertRows(tableName string, newRows []types.RowData) ([]int64, error) {
	ids := make([]int64, len(newRows))
	if len(newRows) == 0 {
		return ids, apperrors.InsertWithNoRows
	}

	var cols []string         // column names for `INSERT INTO (col1, col2...)`
	var args []any            // values to pass to QueryRow
	var placeholders []string // incrementing placeholders e.g. `VALUES (($1, $2), ($3, $4)...)`
	var p int                 // tracks the value of the placeholder

	for k := range newRows[0] {
		cols = append(cols, k)
	}
	// Sort cols so we insert values alphabetically
	slices.Sort(cols)

	for _, newRow := range newRows {
		rowPlaceholders := make([]string, len(cols))
		colIdx := 0
		// Append values in corresponding order of cols
		for _, col := range cols {
			val := newRow[col]
			args = append(args, val)
			rowPlaceholders[colIdx] = fmt.Sprintf("$%d", p+1)
			p++
			colIdx++
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ",")))
	}

	// Build create query
	createStmnt := fmt.Sprintf("INSERT INTO %s (", tableName) +
		strings.Join(cols, ", ") +
		") VALUES " +
		strings.Join(placeholders, ",") +
		" RETURNING id"

	log.Printf("Exec query\n\t%s\nValues: %v\n", createStmnt, args)

	// Execute insert query
	rows, err := r.DB.Query(createStmnt, args...)
	if err != nil {
		return ids, err
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		rows.Scan(&ids[i])
		i++
	}

	return ids, nil
}

func (r *Repository) UpdateRowsByRSQL(tableName string, conditions []rsql.Condition, updatedRow *types.RowData) (int64, error) {
	var assignments []string
	var values []any
	var assignmentVals []any
	var placeholder int
	for k, v := range *updatedRow {
		assignments = append(assignments, fmt.Sprintf("%s = $%d", k, placeholder+1))
		assignmentVals = append(assignmentVals, v)
		placeholder++
	}
	conditional, conditionalVals, err := buildWhereConditions(conditions, placeholder)
	if err != nil {
		return -1, err
	}

	values = slices.Concat(assignmentVals, conditionalVals)

	// Build update query
	updateStmnt := fmt.Sprintf(
		"UPDATE %s SET %s %s",
		tableName,
		strings.Join(assignments, ", "),
		conditional,
	)

	log.Printf("Exec query\n\t%s\nValues: %v", updateStmnt, values)

	// Execute update query
	result, err := r.DB.Exec(updateStmnt, values...)
	if err != nil {
		return -1, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rowsAffected == 0 {
		return -1, apperrors.NewUpdateNoMatchingConditionsErr(tableName, conditions)
	}
	return rowsAffected, nil
}

// DeleteRowByID removes a row from a table by id
func (r *Repository) DeleteRowByID(tableName string, id int64) (int64, error) {
	deleteStmt := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	return r.execDeleteQuery(deleteStmt, []any{id})
}

// DeleteRowsByRSQL removes any rows matching the Condition in the Query
func (r *Repository) DeleteRowsByRSQL(tableName string, conditions []rsql.Condition) (int64, error) {
	// Do not exec delete with empty query
	if len(conditions) == 0 {
		return -1, apperrors.DeleteWithNoConditions
	}
	conditional, values, err := buildWhereConditions(conditions, 0)
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
