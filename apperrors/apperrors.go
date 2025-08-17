package apperrors

import (
	"errors"
	"fmt"
)

var DeleteWithNoFilters = errors.New("Will not DELETE with no filters")

func NewDeleteInvalidIDErr(tableName string, id int64) error {
	return fmt.Errorf(
		"cannot delete row %d in table %s, does not exist",
		id, tableName,
	)
}

func NewUpdateInvalidIDErr(tableName string, id int64) error {
	return fmt.Errorf(
		"cannot update row %d in table %s, does not exist",
		id, tableName,
	)
}
