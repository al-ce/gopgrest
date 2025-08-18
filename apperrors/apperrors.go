package apperrors

import (
	"errors"
	"fmt"

	"gopgrest/rsql"
)

var (
	InsertWithNoRows = errors.New("Cannot insert with no rows")
	DeleteWithNoFilters = errors.New("Will not DELETE with no filters")
)

func NewDeleteInvalidIDErr(tableName string, id int64) error {
	return fmt.Errorf(
		"cannot delete, no row in table %s with id %d",
		tableName, id,
	)
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
