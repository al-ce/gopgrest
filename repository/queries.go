package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"ftrack/types"
)

// ListRows gets rows from a table with optional filter params
func (r *Repository) ListRows(tableName string, qf types.QueryFilters) (*sql.Rows, error) {
	// Build list query with optional conditional filters
	conditional, values, err := buildConditionalClause(qf)
	if err != nil {
		return nil, err
	}
	listStmt := "select * from " + tableName + conditional

	// Execute list query
	rows, err := r.db.Query(listStmt, values...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetRowByID gets a row from a table by id
func (r *Repository) GetRowByID(tableName string, id int) *sql.Row {
	return r.db.QueryRow(
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", tableName),
		id,
	)
}

// InsertRow inserts a new row into a specified table
func (r *Repository) InsertRow(tableName string, newRow *types.RowDataMap) (result InsertResult) {
	// Create cols/values/placeholders slices in consistent order
	var cols []string
	var values []any
	var placeholders []string
	var i int
	for k, v := range *newRow {
		cols = append(cols, k)
		values = append(values, v)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		i++
	}

	// Build create query
	createStmnt := fmt.Sprintf("insert into %s (", tableName) +
		strings.Join(cols, ", ") +
		") values (" +
		strings.Join(placeholders, ",") +
		")" +
		"RETURNING id"

	// Execute insert query
	row := r.db.QueryRow(createStmnt, values...)
	result.Error = row.Scan(&result.ID)
	return
}

// UpdateRowCol updates a field in a table row by id
func (r *Repository) UpdateRowCol(tableName, id, field string, value any) error {
	// Build update query
	updateStmt := fmt.Sprintf(
		"update %s set %s = $1 where id = $2",
		tableName, field,
	)

	if _, err := r.db.Exec(updateStmt, value, id); err != nil {
		return err
	}
	return nil
}

// DeleteRow removes a row from a table by id
func (r *Repository) DeleteRow(tableName, id string) error {
	deleteStmt := fmt.Sprintf("delete from %s where id = $1", tableName)

	// Execute delete query
	result, err := r.db.Exec(deleteStmt, id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
