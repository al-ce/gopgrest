package repository

import (
	"database/sql"
	"fmt"
	"log"
)

// TableColumn represents a column in a table
type TableColumn struct {
	Name     string
	Datatype sql.ColumnType
}

// ColumnMap is a map of column names in a row and their type
type ColumnMap map[string]sql.ColumnType

// Table represents a table in the database
// The Columns slice preserves the column order.
// The ColumnMap is used for fast lookup to check if a column exists
type Table struct {
	Name      string
	Columns   []TableColumn
	ColumnMap ColumnMap
}

// QueryExecutor is an interface that can be satisfied by both *sql.DB and *sql.Tx
type QueryExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// Repository handles database transactions
type Repository struct {
	db     QueryExecutor
	tables []Table
}

// InsertResult contains information from the result of an insert query
type InsertResult struct {
	ID           int64
	Error          error
}

// NewRepository returns a new Repository
func NewRepository(db QueryExecutor, tables []Table) Repository {
	return Repository{
		db:     db,
		tables: tables,
	}
}

// NewTable returns a new Table struct if tableName is a valid table in the
// database
func NewTable(db QueryExecutor, tableName string) (*Table, error) {
	// Get a dummy row of the table with `limit 0`
	query := fmt.Sprintf(`select * from %s limit 0;`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return &Table{}, err
	}
	defer rows.Close()

	var tableColumns []TableColumn
	columnMap := ColumnMap{}

	// Get column names
	colnames, err := rows.Columns()
	if err != nil {
		return &Table{}, err
	}

	// Get column types
	coltypes, err := rows.ColumnTypes()
	if err != nil {
		return &Table{}, err
	}

	// Build slice of TableColumn structs
	for i, name := range colnames {
		tableColumns = append(
			tableColumns,
			TableColumn{
				name,
				*coltypes[i],
			},
		)
		columnMap[name] = *coltypes[i]
	}

	return &Table{
		tableName,
		tableColumns,
		columnMap,
	}, nil
}

// GetPublicTables gets the public tables in the database and builds a slice of
// Table structs to assign to the table field of the Repository
func GetPublicTables(db QueryExecutor) ([]Table, error) {
	// Get table names with query
	rows, err := db.Query(
		`select tablename from pg_catalog.pg_tables where schemaname='public'`,
	)
	if err != nil {
		return []Table{}, err
	}
	defer rows.Close()

	// Build the slice of tables for the Repository
	var tables []Table
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return []Table{}, err
		}

		// Make a new table
		newTable, err := NewTable(db, tableName)
		if err != nil {
			return []Table{}, err
		}
		tables = append(tables, *newTable)
	}

	// Log the tables
	log.Println("Found tables in database:")
	for _, table := range tables {
		log.Printf("\t%s : %d cols", table.Name, len(table.Columns))
		for _, col := range table.Columns {
			log.Printf("\t\t%-15s\t%s", col.Name, col.Datatype.ScanType())
		}
	}
	log.Println()

	return tables, nil
}
