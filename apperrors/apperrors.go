package apperrors

import (
	"errors"
	"fmt"

	"gopgrest/rsql"
)

var (
	InsertWithNoRows         = errors.New("Cannot insert with no rows")
	InsertColsDoNotMatch     = errors.New("Columns in rows to insert do not match")
	InsertValTypesDoNotMatch = errors.New("Value types in rows to insert do not match")
	DeleteWithNoFilters      = errors.New("Will not DELETE with no filters")

	TableDoesNotExist = errors.New("Table does not exist")
	ColDoesNotExist   = errors.New("Column not found in given table")
)

func NewDeleteInvalidIDErr(tableName string, id int64) error {
	return fmt.Errorf("cannot delete, no row in table %s with id %d", tableName, id)
}

func NewUpdateInvalidIDErr(tableName string, id int64) error {
	return fmt.Errorf(
		"cannot update, no row in table %s with id %d",
		tableName, id,
	)
}

func NewUpdateNoMatchingFiltersErr(tableName string, filters []rsql.Filter) error {
	return fmt.Errorf(
		"cannot update, no rows in table %s with filters %v",
		tableName, filters,
	)
}
