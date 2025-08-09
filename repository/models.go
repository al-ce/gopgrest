package repository

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

// TableColumn represents a column in a table
type TableColumn struct {
	Name string
	Type reflect.Type
}

// ColumnMap is a map of column names in a row and their type
type ColumnMap map[string]reflect.Type

// Table represents a table in the database
// The Columns slice preserves the column order.
// The ColumnMap is used for fast lookup to check if a column exists
type Table struct {
	Name      string
	Columns   []TableColumn
	ColumnMap ColumnMap
}

// ColData is a tuple of a column's name and its database type used for JSON
// marshalling, vs. TableColumn which needs reflect.Type to get a type's size
type ColData struct {
	Name string `json:"col_name"`
	Type string `json:"col_type"`
}

// TablesRepr is a map representing database tables and their cols, used for
// JSON marshalling, vs. an array of Table structs which is used for looking up
// the byte size of column types
type TablesRepr map[string][]ColData

// QueryExecutor is an interface that can be satisfied by both *sql.DB and *sql.Tx
type QueryExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// Repository handles database transactions
type Repository struct {
	DB         QueryExecutor
	Tables     []Table
	TablesRepr TablesRepr
}

// ExecResult contains information from the result of an insert query
type ExecResult struct {
	ID    int64
	Error error
}

// NewRepository returns a new Repository
func NewRepository(db QueryExecutor, tables []Table) Repository {
	return Repository{
		DB:         db,
		Tables:     tables,
		TablesRepr: NewTablesRepr(tables),
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
				coltypes[i].ScanType(),
			},
		)
		columnMap[name] = coltypes[i].ScanType()
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
		`SELECT tablename FROM Pg_catalog.pg_tables WHERE schemaname='public'`,
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

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows for db tables: %v", err)
		return nil, err
	}

	// Log the tables
	log.Println("Found tables in database:")
	for _, table := range tables {
		log.Printf("\t%s : %d cols", table.Name, len(table.Columns))
		for _, col := range table.Columns {
			log.Printf("\t\t%-15s\t%s", col.Name, col.Type)
		}
	}
	log.Println()

	return tables, nil
}

func NewTablesRepr(tables []Table) TablesRepr {
	tablesRep := TablesRepr{}
	for _, table := range tables {
		columns := []ColData{}
		for _, col := range table.Columns {
			col := ColData{Name: col.Name, Type: col.Type.Name()}
			columns = append(columns, col)
		}
		tablesRep[table.Name] = columns
	}
	return tablesRep
}
