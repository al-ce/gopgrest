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

// Table represents a table in the database
// The Columns slice preserves the column order.
// The ColumnMap is used for fast lookup to check if a column exists
type Table struct {
	Name      string
	Columns   []TableColumn
	ColumnMap map[string]struct{}
}

// Repository handles database transactions
type Repository struct {
	db     *sql.DB
	tables []Table
}

// NewRepository returns a new Repository
func NewRepository(db *sql.DB) Repository {
	tables, err := getPublicTables(db)
	if err != nil {
		panic(err)
	}
	return Repository{
		db:     db,
		tables: tables,
	}
}

// NewTable returns a new Table struct if tableName is a valid table in the
// database
func NewTable(db *sql.DB, tableName string) (*Table, error) {
	// Get a dummy row of the table with `limit 0`
	query := fmt.Sprintf(`select * from %s limit 0;`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return &Table{}, err
	}
	defer rows.Close()

	var tableColumns []TableColumn
	columnMap := make(map[string]struct{})

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
		columnMap[name] = struct{}{}
	}

	return &Table{
		tableName,
		tableColumns,
		columnMap,
	}, nil
}

// getPublicTables gets the public tables in the database and builds a slice of
// Table structs to assign to the table field of the Repository
func getPublicTables(db *sql.DB) ([]Table, error) {
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
		log.Printf("\t%s\n", table.Name)
		for _, col := range table.Columns {
			log.Printf("\t\t%-15s\t%s", col.Name, col.Datatype.ScanType())
		}
	}
	log.Println()

	return tables, nil
}
