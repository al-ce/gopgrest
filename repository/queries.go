package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"ftrack/types"
)

// ListRows gets rows from a table with optional filter params
func (r *Repository) ListRows(tableName string, qf types.QueryFilter) (*sql.Rows, error) {
	// Build list query with optional conditional filters
	conditional, conditionalLog, values, err := buildConditionalClause(qf)
	if err != nil {
		return nil, err
	}

	log.Printf("SELECT * FROM %s%s", tableName, conditionalLog)

	listStmt := "SELECT * FROM " + tableName + conditional

	// Execute list query
	rows, err := r.DB.Query(listStmt, values...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetRowByID gets a row from a table by id
func (r *Repository) GetRowByID(tableName, id string) *sql.Row {
	log.Printf(
		"Exec query\n\tSELECT * FROM %s WHERE id = %s",
		tableName, id,
	)

	return r.DB.QueryRow(
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", tableName),
		id,
	)
}

// InsertRow inserts a new row into a specified table
func (r *Repository) InsertRow(tableName string, newRow *types.RowData) (result InsertResult) {
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

	log.Printf(
		"Exec query\n\tINSERT INTO %s (%s) values (%s) RETURNING id",
		tableName, strings.Join(cols, ", "), strings.Join(valuesLog, ", "),
	)

	// Build create query
	createStmnt := fmt.Sprintf("INSERT INTO %s (", tableName) +
		strings.Join(cols, ", ") +
		") values (" +
		strings.Join(placeholders, ",") +
		")" +
		"RETURNING id"

	// Execute insert query
	row := r.DB.QueryRow(createStmnt, values...)
	result.Error = row.Scan(&result.ID)
	return
}

// UpdateRowCol updates a column in a table row by id
func (r *Repository) UpdateRowCol(tableName, id, col string, value any) error {
	log.Printf(
		"Exec query\n\tUPDATE %s SET %s = %s WHERE id = %s\n",
		tableName, col, value, id,
	)

	// Build update query
	updateStmt := fmt.Sprintf(
		"UPDATE %s SET %s = $1 WHERE id = $2",
		tableName, col,
	)

	if _, err := r.DB.Exec(updateStmt, value, id); err != nil {
		return err
	}

	return nil
}

// DeleteRow removes a row from a table by id
func (r *Repository) DeleteRow(tableName, id string) error {
	deleteStmt := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	log.Printf(
		"Exec query \n\tDELETE FROM %s WHERE id = %s\n",
		tableName, id,
	)

	// Execute delete query
	result, err := r.DB.Exec(deleteStmt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf(
			"row %s in table %s does not exist, did not attempt delete",
			id, tableName,
		)
	}

	return nil
}
