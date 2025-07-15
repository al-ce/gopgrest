package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

// ListRows gets rows from a table with optional filter params
func (r *Repository) ListRows(table string, params map[string][]string) (*sql.Rows, error) {
	// Build list query with optional conditional filters
	conditional, values := buildConditionalClause(params)
	listStmt := "select * from " + table + conditional

	// Execute list query
	rows, err := r.db.Query(listStmt, values...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// InsertRow inserts a new row into a specified table
func (r *Repository) InsertRow(table string, newRow *map[string]any) error {
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
	createStmnt := fmt.Sprintf("insert into %s (", table) +
		strings.Join(cols, ", ") +
		") values (" +
		strings.Join(placeholders, ",") +
		")"

	// Execute create query
	result, err := r.db.Exec(createStmnt, values...)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// UpdateRowCol updates a field in a table row by id
func (r *Repository) UpdateRowCol(table, id, field string, value any) error {
	// Build update query
	updateStmt := fmt.Sprintf(
		"update %s set %s = $1 where id = $2",
		table, field,
	)

	if _, err := r.db.Exec(updateStmt, value, id); err != nil {
		return err
	}
	return nil
}

// DeleteRow removes a row from a table by id
func (r *Repository) DeleteRow(table, id string) error {
	deleteStmt := fmt.Sprintf("delete from %s where id = $1", table)

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
